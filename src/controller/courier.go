package controller

import (
	"database/sql"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"yandex-team.ru/bstask/models"
)

type CreateCouriersRequest struct {
	Couriers []models.CreateCourierDto `json:"couriers" validate:"required,dive,required"`
}

type CreateCouriersResponse struct {
	Couriers []models.CourierDto `json:"couriers"`
}

type GetCouriersResponse struct {
	Couriers []models.CourierDto `json:"couriers"`
	Limit    int32               `json:"limit"`
	Offset   int32               `json:"offset"`
}

// e.GET("/couriers", getCouriers)
func (c *Controller) getCouriers(ctx echo.Context) error {
	var params struct {
		Limit  *int32 `query:"limit" validate:"omitempty,gte=1"`
		Offset *int32 `query:"offset" validate:"omitempty,gte=0"`
	}

	if err := ctx.Bind(&params); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&params); err != nil {
		return badRequest(ctx, err)
	}

	limit := getOptional(params.Limit, 1)
	offset := getOptional(params.Offset, 0)

	courierList, err := c.Service.GetCouriers(limit, offset)
	if err != nil {
		return badRequest(ctx, err)
	}

	return ctx.JSON(
		http.StatusOK,
		GetCouriersResponse{
			Couriers: courierList,
			Limit:    limit,
			Offset:   offset})
}

// e.POST("/couriers", createCouriers)
func (c *Controller) createCouriers(ctx echo.Context) error {
	var courierRequest CreateCouriersRequest

	if err := ctx.Bind(&courierRequest); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&courierRequest); err != nil {
		return badRequest(ctx, err)
	}

	courierList, err := c.Service.CreateCouriers(courierRequest.Couriers)

	if err != nil {
		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, CreateCouriersResponse{Couriers: courierList})
}

// e.GET("/couriers/:courier_id", getCourierById)
func (c *Controller) getCourierById(ctx echo.Context) error {
	var params struct {
		CourierId int64 `param:"courier_id" validate:"gte=1"`
	}

	if err := ctx.Bind(&params); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&params); err != nil {
		return badRequest(ctx, err)
	}

	courier, err := c.Service.GetCourierById(params.CourierId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(ctx, err)
		}

		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, courier)
}

// e.GET("/couriers/meta-info/:courier_id", getCourierMetaInfo)
func (c *Controller) getCourierMetaInfo(ctx echo.Context) error {
	var params struct {
		CourierId int64  `param:"courier_id" validate:"gte=1"`
		StartDate string `query:"startDate" validate:"isValidDate"`
		EndDate   string `query:"endDate" validate:"isValidDate"`
	}

	if err := ctx.Bind(&params); err != nil {
		return badRequest(ctx, err)
	}

	if err := ctx.Validate(&params); err != nil {
		return badRequest(ctx, err)
	}

	startTime, _ := time.Parse(models.Date, params.StartDate)
	endTime, _ := time.Parse(models.Date, params.EndDate)

	if !endTime.After(startTime) {
		return badRequest(ctx, errors.New("end date should be > start date"))
	}

	courierMetaInfo, err := c.Service.GetCourierMetaInfoById(params.CourierId, params.StartDate, params.EndDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return notFound(ctx, err)
		}

		return badRequest(ctx, err)
	}

	return ctx.JSON(http.StatusOK, courierMetaInfo)
}
