package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	Cache    *redis.Client
	Ctx      context.Context
	Messages <-chan *redis.Message
}

func NewCacheClient() *CacheClient {
	client := InitRedisClient()
	ctx, _ := context.WithTimeout(context.TODO(), time.Millisecond*10)

	cacheClient := &CacheClient{
		Cache: client,
		Ctx:   ctx,
	}
	cacheClient.Subscribe()
	return cacheClient
}

func (cache *CacheClient) Subscribe() {
	pubsub := cache.Cache.Subscribe(cache.Ctx, "chatMessage")
	cache.Messages = pubsub.Channel()
}
