package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func badRequest(ctx echo.Context, err error) error {
	ctx.Logger().Error(err)
	return ctx.JSON(http.StatusBadRequest, struct{}{})
}

func notFound(ctx echo.Context, err error) error {
	ctx.Logger().Error(err)
	return ctx.JSON(http.StatusNotFound, struct{}{})
}

func getOptional[T any](p *T, defaultVal T) T {
	val := defaultVal
	if p != nil {
		val = *p
	}

	return val
}
