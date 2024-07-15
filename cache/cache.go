package cache

import (
	"context"
	"egame_core/config"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

// Initialize initializes the Redis client
func Initialize() {
	config := config.GetConfig()
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host + ":" + config.Redis.Port, // Redis 服務地址
		Password: config.Redis.Password,                       // 沒有密碼設置
		DB:       0,                                           // 使用默認的 DB
	})

	// 測試連接
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
}

func GetClient() *redis.Client {
	if rdb == nil {
		Initialize()
	}
	return rdb
}

// Close closes the Redis client
func Close() {
	err := rdb.Close()
	if err != nil {
		log.Fatalf("failed to close Redis client: %v", err)
	}
}
