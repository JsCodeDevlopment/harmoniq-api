package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Println("Redis not connected. Ensure redis is running or ignore if not needed full-app:", err)
		return
	}
	log.Println("Redis connected successfully!")
}
