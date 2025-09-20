package routes

import (
	"github.com/labstack/echo/v4"
	"microservice/internal/modules/tenant"
)

func Tenant(e *echo.Group, h tenant.ITenantHttpHandler) {
	r := e.Group("/tenant")
	r.POST("/create", h.Create)
	r.GET("/:uuid", h.Details)
	r.GET("/list", h.List)
}
