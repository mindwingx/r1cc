package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "microservice/docs" // the custom path of the generated swagger files
	"microservice/internal/server/http/routes"
	"strings"
)

func (s *Server) setRoutes() {
	//general routes
	s.client.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	s.client.GET("/handshake", s.health.Handshake)
	s.client.GET("/public/swagger/*", routes.Swagger(s.swagger), s.middleware.SwagAuth(s.swagger))

	api := s.client.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			routes.Tenant(v1, s.tenant)
			routes.Credit(v1, s.credit)
			routes.Message(v1, s.message)
		}
	}
}

func (s *Server) routesStdout() {
	if routesCount := len(s.client.Routes()); routesCount > 0 {
		rs := strings.Builder{}
		rs.WriteString("\n")
		defer func() { rs.WriteString("\n"); fmt.Print(rs.String()) }()

		for _, r := range s.client.Routes() {
			rs.WriteString(fmt.Sprintf("[%s] %s - %s\n", r.Method, r.Path, r.Name))
		}
	}
}
