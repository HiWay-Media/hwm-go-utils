package generic

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetEndpoints[T any](router fiber.Router, database *gorm.DB, logger *zap.SugaredLogger) (IStore[T], IService[T], IHandler[T]) {
	store := NewStore[T](database)
	service := NewService[T](store, logger)
	handler := NewHandler[T](service, logger)
	setRoutes(router, handler)
	return store, service, handler
}

func setRoutes[T any](router fiber.Router, handler IHandler[T]) {
	router.Get("/:id", handler.Get)
	router.Get("/", handler.List)
	router.Post("/", handler.Create)
	router.Delete("/:id", handler.Delete)
}
