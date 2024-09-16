package infrastructure

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedis(ctx context.Context) (*RedisClient, error) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST and REDIS_PORT environment variables are required")
	}
	log.Printf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	ping, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("redis connection successful:", ping)
	return &RedisClient{Client: rdb}, nil
}
