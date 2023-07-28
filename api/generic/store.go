package generic

import "gorm.io/gorm"

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

func (s *Store[T]) List(start, limit int) ([]T, error) {
	var records []T

	result := s.db.Offset(start).Limit(limit).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return records, nil
}

func (s *Store[T]) Get(id any) (*T, error) {
	var t T
	result := s.db.First(&t, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &t, nil
}

func (s *Store[T]) Create(obj *T) error {
	result := s.db.Create(&obj)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Store[T]) Delete(id any) error {
	var t T
	result := s.db.Delete(&t, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
