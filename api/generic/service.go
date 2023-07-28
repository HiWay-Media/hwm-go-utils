package generic

import (
	"go.uber.org/zap"
)

type IService[TK any, T any] interface {
	Get(id TK) (*T, error)
	Create(obj *T) error
	Delete(id TK) error
	List(start, limit int) ([]T, error)
}

type Service[TK any, T any] struct {
	Store  IStore[TK, T]
	Logger *zap.SugaredLogger
}

func NewService[TK any, T any](store IStore[TK, T], logger *zap.SugaredLogger) IService[TK, T] {
	return &Service[TK, T]{Logger: logger, Store: store}
}

func (s *Service[TK, T]) List(start, limit int) ([]T, error) {
	return s.Store.List(start, limit)
}

func (s *Service[TK, T]) Get(id TK) (*T, error) {
	return s.Store.Get(id)
}

func (s *Service[TK, T]) Create(obj *T) error {
	return s.Store.Create(obj)
}

func (s *Service[TK, T]) Delete(id TK) error {
	return s.Store.Delete(id)
}
