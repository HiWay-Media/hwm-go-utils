package generic

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetEndpoints[TK any, T any](router fiber.Router, database *gorm.DB, logger *zap.SugaredLogger) (IStore[TK, T], IService[TK, T], IHandler[TK, T]) {
	store := NewStore[TK, T](database)
	service := NewService[TK, T](store, logger)
	handler := NewHandler[TK, T](service, logger)
	setRoutes(router, handler)
	return store, service, handler
}

func setRoutes[TK any, T any](router fiber.Router, handler IHandler[TK, T]) {
	router.Get("/:id", handler.Get)
	router.Get("/", handler.List)
	router.Post("/", handler.Create)
	router.Delete("/:id", handler.Delete)
}
