package routes

import (
	"github.com/labstack/echo/v4"
	swagger "github.com/swaggo/echo-swagger"
	"microservice/config"
	"microservice/docs"
)

func Swagger(swg *config.Swagger) echo.HandlerFunc {
	if swg == nil || swg.Enable == false {
		return func(c echo.Context) error { return nil }
	}

	docs.SwaggerInfo.Schemes = []string{swg.Schemes}
	docs.SwaggerInfo.Host = swg.Host
	docs.SwaggerInfo.Title = swg.Title
	docs.SwaggerInfo.Description = swg.Description
	docs.SwaggerInfo.Version = swg.Version

	return swagger.EchoWrapHandler(func(config *swagger.Config) {
		config.DocExpansion = "none"
		config.DeepLinking = true
		config.SyntaxHighlight = true
	})
}
