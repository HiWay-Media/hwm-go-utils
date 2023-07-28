package generic

import (
	"go.uber.org/zap"
)

type IService[T any] interface {
	Get(id any) (*T, error)
	Create(obj *T) error
	Delete(id any) error
	List(start, limit int) ([]T, error)
}

type Service[T any] struct {
	Store  IStore[T]
	Logger *zap.SugaredLogger
}

func NewService[T any](store IStore[T], logger *zap.SugaredLogger) IService[T] {
	return &Service[T]{Logger: logger, Store: store}
}

func (s *Service[T]) List(start, limit int) ([]T, error) {
	return s.Store.List(start, limit)
}

func (s *Service[T]) Get(id any) (*T, error) {
	return s.Store.Get(id)
}

func (s *Service[T]) Create(obj *T) error {
	return s.Store.Create(obj)
}

func (s *Service[T]) Delete(id any) error {
	return s.Store.Delete(id)
}
