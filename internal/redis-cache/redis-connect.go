package cache

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedisClient() *redis.Client {
	var ctx = context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	redisClient.Ping(ctx)
	return redisClient
}

func GetUserWS(ctx context.Context, userId string) {
	val, err := redisClient.Get(ctx, userId).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Set value as ", val)

}

func SetUserWS(ctx context.Context, conn net.Conn, userId string) {
	_, err := redisClient.Set(ctx, userId, conn, 1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Set value Successfully.")
}
