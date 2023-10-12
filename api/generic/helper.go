package generic

import (
	"github.com/Paxx-RnD/go-helper/helpers/generic_helper"
	"github.com/gofiber/fiber/v2"
)

func FromQuery[T any](ctx *fiber.Ctx, param string, defaultValue T) (*T, error) {
	q := ctx.Query(param, "")
	v, err := generic_helper.ConvertFromString[T](q)
	if defaultValue != nil && v == nil {
		return &defaultValue, err
	}
	return v, err
}

func FromParam[T any](ctx *fiber.Ctx, param string, defaultValue T) (*T, error) {
	p := ctx.Params(param, "")
	v, err := generic_helper.ConvertFromString[T](p)
	if defaultValue != nil && v == nil {
		return &defaultValue, err
	}
	return v, err
}
