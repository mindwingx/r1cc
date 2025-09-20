package middleware

import (
	"github.com/labstack/echo/v4"
	"microservice/config"
)

type IMiddleware interface {
	Service() *config.Service
	SwagAuth(swg *config.Swagger) echo.MiddlewareFunc
	RequestCounter(next echo.HandlerFunc) echo.HandlerFunc
	RequestDuration(next echo.HandlerFunc) echo.HandlerFunc
	RequestProcess(next echo.HandlerFunc) echo.HandlerFunc
}
