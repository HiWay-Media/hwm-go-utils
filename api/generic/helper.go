package generic

import (
	"github.com/Paxx-RnD/go-helper/helpers/generic_helper"
	"github.com/gofiber/fiber/v2"
)

func FromQuery[T any](ctx *fiber.Ctx, param string, out *T) error {
	q := ctx.Query(param, "")
	return generic_helper.ConvertFromString[T](q, out)
}

func FromParam[T any](ctx *fiber.Ctx, param string, out *T) error {
	p := ctx.Params(param, "")
	return generic_helper.ConvertFromString[T](p, out)
}
