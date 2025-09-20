package routes

import (
	"github.com/labstack/echo/v4"
	"microservice/internal/modules/credit"
)

func Credit(e *echo.Group, h credit.ICreditHttpHandler) {
	r := e.Group("/credit")
	r.POST("/increase", h.IncreaseCredit)
	r.GET("/transactions", h.TransactionsList)
}
