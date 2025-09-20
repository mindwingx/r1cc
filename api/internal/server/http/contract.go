package http

import "go.uber.org/fx"

type IHttpServer interface {
	Init()
	Fx(lc fx.Lifecycle, sfx ServerFx)
}
