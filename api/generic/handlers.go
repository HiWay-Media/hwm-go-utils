package generic

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"sync"

	"github.com/HiWay-Media/hwm-go-utils/api/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type IHandler[T any] interface {
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

type Handler[T any] struct {
	Service IService[T]
	Logger  *zap.SugaredLogger
}

func NewHandler[T any](service IService[T], logger *zap.SugaredLogger) IHandler[T] {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	return &Handler[T]{Service: service, Logger: logger}
}

// errorResponse maps a service error to a response without leaking database
// details to the client; unexpected errors are logged instead.
func (s *Handler[T]) errorResponse(c *fiber.Ctx, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(http.StatusNotFound).JSON(models.ApiDefaultError("not found"))
	}
	s.Logger.Errorf("%s %s: %v", c.Method(), c.Path(), err)
	return c.Status(http.StatusInternalServerError).JSON(models.ApiDefaultError("internal error"))
}

func (s *Handler[T]) List(c *fiber.Ctx) error {
	start, err := strconv.Atoi(c.Query("start", "0"))
	if err != nil || start < 0 {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("start invalid"))
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err != nil || limit < 0 {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("limit invalid"))
	}

	r, err := s.Service.List(start, limit)
	if err != nil {
		return s.errorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(r))
}

func (s *Handler[T]) Create(c *fiber.Ctx) error {
	var requestBody T

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError(err.Error()))
	}

	if err := validator.Validate(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError(err.Error()))
	}

	if err := clearServerOwnedFields(&requestBody); err != nil {
		return s.errorResponse(c, err)
	}

	if err := s.Service.Create(&requestBody); err != nil {
		return s.errorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(requestBody))
}

var schemaCache sync.Map

// clearServerOwnedFields zeroes the primary key and the relationship fields of
// a request body, so a client cannot choose the id of a new row or make GORM
// create/link associated rows through a generic Create endpoint.
func clearServerOwnedFields(obj any) error {
	sch, err := schema.Parse(obj, &schemaCache, schema.NamingStrategy{})
	if err != nil {
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(obj))
	fields := append([]*schema.Field{}, sch.PrimaryFields...)
	for _, rel := range sch.Relationships.Relations {
		fields = append(fields, rel.Field)
	}
	for _, f := range fields {
		if err := f.Set(context.Background(), rv, reflect.Zero(f.FieldType).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func parseID(c *fiber.Ctx) (any, bool) {
	id := c.Params("id")
	if id == "" {
		return nil, false
	}
	if n, err := strconv.Atoi(id); err == nil {
		return n, true
	}
	return id, true
}

func (s *Handler[T]) Get(c *fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("id invalid"))
	}

	r, err := s.Service.Get(id)
	if err != nil {
		return s.errorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(r))
}

func (s *Handler[T]) Delete(c *fiber.Ctx) error {
	id, ok := parseID(c)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("id invalid"))
	}

	if err := s.Service.Delete(id); err != nil {
		return s.errorResponse(c, err)
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefault())
}
