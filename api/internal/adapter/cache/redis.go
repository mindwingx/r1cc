package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"microservice/config"
	"microservice/internal/adapter/registry"
	"microservice/pkg/utils"
	"time"
)

type cache struct {
	redis.Client
	config config.Cache
}

func New(registry registry.IRegistry) ICache {
	c := new(cache)
	if err := registry.Parse(&c.config); err != nil {
		utils.PrintStd(utils.StdPanic, "cache", "config parse err: %s", err)
	}

	return c
}

func (c *cache) Init() {
	var err error

	c.Client = *redis.NewClient(&redis.Options{
		DB:       c.config.DB,
		Addr:     fmt.Sprintf("%s:%s", c.config.Host, c.config.Port),
		Username: c.config.Username,
		Password: c.config.Password,
	})

	if err = c.Client.Ping(context.Background()).Err(); err != nil {
		zap.L().Error(err.Error())
		utils.PrintStd(utils.StdPanic, "cache", "init err: %s", err)
	}
}

// C redis client instance
func (c *cache) C() (r *redis.Client) {
	r = &c.Client
	return
}

// Set meth a new key,value
func (c *cache) Set(ctx context.Context, key string, value interface{}, duration time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return c.Client.Set(ctx, key, p, duration).Err()
}

// Get meth, get value with key
func (c *cache) Get(ctx context.Context, key string, dest interface{}) error {
	p, err := c.Client.Get(ctx, key).Result()

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(p), &dest)
}

// Del for delete keys in redis
func (c *cache) Del(ctx context.Context, key ...string) error {
	_, err := c.Client.Del(ctx, key...).Result()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	return nil
}

func (c *cache) Fx(lc fx.Lifecycle) ICache {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "cache", "initiated")
			return
		},
		OnStop: func(ctx context.Context) (err error) {
			utils.PrintStd(utils.StdLog, "cache", "stopping...")

			if err = c.Client.Close(); err != nil {
				utils.PrintStd(utils.StdLog, "cache", "connection close", err)
				return
			}

			utils.PrintStd(utils.StdLog, "cache", "stopped")
			return
		},
	})

	return c
}
