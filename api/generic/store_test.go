package generic

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type item struct {
	ID   uint
	Name string
}

// dryRunStore returns a Store backed by a dry-run DB and a func yielding the last SQL it built.
func dryRunStore(t *testing.T) (IStore[item], func() string) {
	t.Helper()
	db, err := gorm.Open(mysql.New(mysql.Config{DSN: "u:p@tcp(127.0.0.1:1)/x", SkipInitializeWithVersion: true}),
		&gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Discard})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}

	var last string
	capture := func(tx *gorm.DB) {
		last = strings.TrimSpace(tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...))
	}
	_ = db.Callback().Query().After("gorm:query").Register("test:capture", capture)
	_ = db.Callback().Delete().After("gorm:delete").Register("test:capture", capture)

	return NewStore[item](db), func() string { return last }
}

func TestStoreBindsIdsAsParameters(t *testing.T) {
	store, lastSQL := dryRunStore(t)

	cases := []struct {
		name string
		run  func()
		want string
	}{
		{
			name: "get numeric",
			run:  func() { _, _ = store.Get(5) },
			want: "SELECT * FROM `items` WHERE `items`.`id` = 5 ORDER BY `items`.`id` LIMIT 1",
		},
		{
			name: "get injection",
			run:  func() { _, _ = store.Get("1 OR 1=1") },
			want: "SELECT * FROM `items` WHERE `items`.`id` = '1 OR 1=1' ORDER BY `items`.`id` LIMIT 1",
		},
		{
			name: "delete injection",
			run:  func() { _ = store.Delete("1 OR 1=1") },
			want: "DELETE FROM `items` WHERE `items`.`id` = '1 OR 1=1'",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.run()
			if got := lastSQL(); got != c.want {
				t.Errorf("got  %s\nwant %s", got, c.want)
			}
		})
	}
}

func TestStoreDeleteMissingRowIsNotFound(t *testing.T) {
	store, _ := dryRunStore(t)
	// a dry run affects no rows
	if err := store.Delete(5); err != gorm.ErrRecordNotFound {
		t.Errorf("got %v, want %v", err, gorm.ErrRecordNotFound)
	}
}

func TestStoreListZeroLimitIsUnbounded(t *testing.T) {
	store, lastSQL := dryRunStore(t)
	_, _ = store.List(0, 0)
	if got, want := lastSQL(), "SELECT * FROM `items`"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
