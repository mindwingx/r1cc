package routes

import (
	"github.com/labstack/echo/v4"
	"microservice/internal/modules/message"
)

func Message(e *echo.Group, h message.IMessageHttpHandler) {
	r := e.Group("/message")
	r.POST("/send", h.Send)
	r.GET("/list", h.List)
}
