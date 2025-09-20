package app

import (
	"microservice/internal/server/http"
)

type Servers struct {
	http http.IHttpServer
}

func NewServers() *Servers { return &Servers{http: &http.Server{}} }

// servers

func (s *Servers) Http() http.IHttpServer        { return s.http }
func (s *Servers) SetHttp(http http.IHttpServer) { s.http = http }

// init

func (a *App) InitServers() {
	a.SetServer(NewServers())
	a.initHttpServer()
}

func (a *App) initHttpServer() {
	server := http.NewServer(a.Config(), a.Client().Registry())
	a.Server().SetHttp(server)
	a.Server().Http().Init()
	a.Span().AddEvent("http server initialized")
}
