package generic

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetEndpoints[T any](group string, router fiber.Router, database *gorm.DB, logger *zap.SugaredLogger) (IStore[T], IService[T], IHandler[T]) {
	store := NewStore[T](database)
	service := NewService[T](store, logger)
	handler := NewHandler[T](service, logger)
	setRoutes(group, router, handler)
	return store, service, handler
}

func setRoutes[T any](group string, router fiber.Router, handler IHandler[T]) {
	g := router.Group("/" + group)
	g.Get("/:id", handler.Get)
	g.Get("/", handler.List)
	g.Post("/", handler.Create)
	g.Delete("/:id", handler.Delete)
}
