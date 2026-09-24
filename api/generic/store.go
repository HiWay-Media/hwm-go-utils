package generic

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IStore[T any] interface {
	Get(id any) (*T, error)
	Create(obj *T) error
	Delete(id any) error
	List(start, limit int) ([]T, error)
}

type Store[T any] struct {
	db *gorm.DB
}

func NewStore[T any](db *gorm.DB) IStore[T] {
	return &Store[T]{db: db}
}

// byPrimaryKey binds id as a parameter against the model's primary key.
// Passing id straight to First/Delete is unsafe: GORM treats a non-numeric
// string there as a raw SQL condition (e.g. "1 OR 1=1").
func byPrimaryKey(id any) clause.Eq {
	return clause.Eq{Column: clause.Column{Table: clause.CurrentTable, Name: clause.PrimaryKey}, Value: id}
}

func (s *Store[T]) List(start, limit int) ([]T, error) {
	var records []T

	// GORM renders Limit(0) as "LIMIT 0"; treat a non-positive limit as "no limit"
	if limit <= 0 {
		limit = -1
	}

	result := s.db.Offset(start).Limit(limit).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (s *Store[T]) Get(id any) (*T, error) {
	var t T
	result := s.db.Where(byPrimaryKey(id)).First(&t)
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}

func (s *Store[T]) Create(obj *T) error {
	result := s.db.Create(obj)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store[T]) Delete(id any) error {
	var t T
	result := s.db.Where(byPrimaryKey(id)).Delete(&t)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
