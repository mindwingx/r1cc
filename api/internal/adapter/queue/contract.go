package queue

import (
	"context"
	"go.uber.org/fx"
)

type IQueue interface {
	Init()
	Produce(ctx context.Context, topic, key string, value []byte) error
	Fx(lc fx.Lifecycle, qfx QFx) IQueue
}
