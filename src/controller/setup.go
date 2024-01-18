package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func (c *Controller) SetupRoutes(e *echo.Echo) {
	e.GET("/ping", ping)

	e.GET("/orders", c.getOrders, newRateLimiter())

	e.POST("/orders", c.createOrders, newRateLimiter())

	e.POST("/orders/complete", c.completeOrders, newRateLimiter())

	e.GET("/couriers", c.getCouriers, newRateLimiter())

	e.POST("/couriers", c.createCouriers, newRateLimiter())

	e.GET("/orders/:order_id", c.getOrderById, newRateLimiter())

	e.GET("/couriers/:courier_id", c.getCourierById, newRateLimiter())

	e.GET("/couriers/meta-info/:courier_id", c.getCourierMetaInfo, newRateLimiter())
}

func newRateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10)))
}
