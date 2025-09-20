package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"microservice/config"
)

func (m *Middleware) SwagAuth(swg *config.Swagger) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Replace with your authentication logic
		if username == swg.Username && password == swg.Password {
			return true, nil
		}

		return false, nil
	})
	// You can add additional logging or enhancements here if needed
}
