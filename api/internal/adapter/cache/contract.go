package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"time"
)

type ICache interface {
	Init()
	C() *redis.Client
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key ...string) error
	Fx(lc fx.Lifecycle) ICache
}
