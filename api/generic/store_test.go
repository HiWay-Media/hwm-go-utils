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

func TestStoreListLimit(t *testing.T) {
	defer func(old int) { MaxListLimit = old }(MaxListLimit)

	cases := []struct {
		name       string
		max, start int
		limit      int
		want       string
	}{
		{"missing limit is capped", 1000, 0, 0, "SELECT * FROM `items` LIMIT 1000"},
		{"large limit is capped", 1000, 0, 5000, "SELECT * FROM `items` LIMIT 1000"},
		{"limit within cap", 1000, 20, 50, "SELECT * FROM `items` LIMIT 50 OFFSET 20"},
		{"cap disabled, missing limit is unbounded", 0, 0, 0, "SELECT * FROM `items`"},
		{"cap disabled, explicit limit", 0, 0, 5000, "SELECT * FROM `items` LIMIT 5000"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			MaxListLimit = c.max
			store, lastSQL := dryRunStore(t)
			_, _ = store.List(c.start, c.limit)
			if got := lastSQL(); got != c.want {
				t.Errorf("got %s, want %s", got, c.want)
			}
		})
	}
}
