package controller

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"yandex-team.ru/bstask/models"
)

type CompleteOrderRequestDto struct {
	CompleteInfo []models.CompleteOrderDto `json:"complete_info" validate:"required,dive,required"`
}

type CreateOrderRequest struct {
	Orders []models.CreateOrderDto `json:"orders" validate:"required,dive,required"`
}

// e.GET("/orders", getOrders)
func (c *Controller) getOrders(ctx echo.Context) error {
	var params struct {
		Limit  *int32 `query:"limit" validate:"omitempty,gte=0"`
		Offset *int32 `query:"offset" validate:"omitempty,gte=0"`
	}

	if err := ctx.Bind(&params); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&params); err != nil {
		return badRequest(ctx, err)
	}

	orderList, err := c.Service.GetOrders(getOptional(params.Limit, 1), getOptional(params.Offset, 0))
	if err != nil {
		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, orderList)
}

// e.POST("/orders", createOrders)
func (c *Controller) createOrders(ctx echo.Context) error {
	var orderRequest CreateOrderRequest

	if err := ctx.Bind(&orderRequest); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&orderRequest); err != nil {
		return badRequest(ctx, err)
	}

	orderList, err := c.Service.CreateOrders(orderRequest.Orders)
	if err != nil {
		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, orderList)
}

// e.POST("/orders/complete", completeOrders)
func (c *Controller) completeOrders(ctx echo.Context) error {
	var completeOrderRequest CompleteOrderRequestDto

	if err := ctx.Bind(&completeOrderRequest); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&completeOrderRequest); err != nil {
		return badRequest(ctx, err)
	}

	orderList, err := c.Service.CompleteOrders(completeOrderRequest.CompleteInfo)
	if err != nil {
		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, orderList)
}

// e.GET("/orders/:order_id", getOrderById)
func (c *Controller) getOrderById(ctx echo.Context) error {
	var params struct {
		OrderId int64 `param:"order_id" validate:"gte=1"`
	}

	if err := ctx.Bind(&params); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&params); err != nil {
		return badRequest(ctx, err)
	}

	order, err := c.Service.GetOrderById(params.OrderId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(ctx, err)
		}

		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, order)
}
