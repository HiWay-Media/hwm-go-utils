package generic

import "gorm.io/gorm"

type IStore[TK any, T any] interface {
	Get(id TK) (*T, error)
	Create(obj *T) error
	Delete(id TK) error
	List(start, limit int) ([]T, error)
}

type Store[TK any, T any] struct {
	db *gorm.DB
}

func NewStore[TK any, T any](db *gorm.DB) IStore[TK, T] {
	return &Store[TK, T]{db: db}
}

func (s *Store[TK, T]) List(start, limit int) ([]T, error) {
	var records []T

	result := s.db.Offset(start).Limit(limit).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (s *Store[TK, T]) Get(id TK) (*T, error) {
	var t T
	result := s.db.First(&t, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}

func (s *Store[TK, T]) Create(obj *T) error {
	result := s.db.Create(&obj)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store[TK, T]) Delete(id TK) error {
	var t T
	result := s.db.Delete(&t, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
