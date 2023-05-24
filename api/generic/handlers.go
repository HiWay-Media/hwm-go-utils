package generic

import (
	"github.com/HiWay-Media/hwm-go-utils/api/models"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type IHandler[TK any, T any] interface {
	Get(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

type Handler[TK any, T any] struct {
	Service       IService[TK, T]
	Configuration Configuration
	Logger        *zap.SugaredLogger
}

func NewHandler[TK any, T any](service IService[TK, T], logger *zap.SugaredLogger, configuration Configuration) IHandler[TK, T] {
	return &Handler[TK, T]{Logger: logger, Configuration: configuration, Service: service}
}

func (s *Handler[TK, T]) List(c *fiber.Ctx) error {
	start, err := strconv.Atoi(c.Query("start", "0"))
	if err == nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("start invalid"))
	}

	limit, err := strconv.Atoi(c.Query("limit", "0"))
	if err == nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("limit invalid"))
	}

	r, err := s.Service.List(start, limit)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(models.ApiDefaultError(err.Error()))
		}
		return c.Status(http.StatusInternalServerError).JSON(models.ApiDefaultError(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(r))
}

func (s *Handler[TK, T]) Create(c *fiber.Ctx) error {
	var requestBody T

	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError(err.Error()))
	}

	if err := validator.Validate(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError(err.Error()))
	}

	err := s.Service.Create(&requestBody)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ApiDefaultError(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(requestBody))
}

func (s *Handler[TK, T]) Get(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("id invalid"))
	}

	var anyId any
	anyId, err := strconv.Atoi(id)
	if err != nil {
		anyId = id
	}

	r, err := s.Service.Get[TK, T](anyId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(models.ApiDefaultError(err.Error()))
		}
		return c.Status(http.StatusInternalServerError).JSON(models.ApiDefaultError(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefaultResponse(r))
}

func (s *Handler[TK, T]) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(models.ApiDefaultError("id invalid"))
	}

	var anyId any
	anyId, err := strconv.Atoi(id)
	if err != nil {
		anyId = id
	}

	err = s.Service.Delete[TK, T](anyId)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(http.StatusNotFound).JSON(models.ApiDefaultError(err.Error()))
		}
		return c.Status(http.StatusInternalServerError).JSON(models.ApiDefaultError(err.Error()))
	}

	return c.Status(http.StatusOK).JSON(models.ApiDefault())
}
