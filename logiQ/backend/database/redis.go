package database

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient 全局Redis客户端实例
var RedisClient *redis.Client

// ConnectRedis 连接到Redis数据库
// 这个函数在应用启动时调用，建立与Redis的连接
func ConnectRedis() {
	// Redis连接地址，从环境变量读取
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set")
	}

	password := os.Getenv("REDIS_PASSWORD")
	redisDB := 0
	dbEnv := os.Getenv("REDIS_DB")
	if dbEnv == "" {
		log.Fatal("REDIS_DB is not set")
	}
	parsed, err := strconv.Atoi(dbEnv)
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value %s: %v", dbEnv, err)
	}
	redisDB = parsed

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password, // 你的Redis密码，如果设置了的话
		DB:       redisDB,  // 默认DB 0
	})

	// 创建上下文，设置5秒超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接是否正常工作
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis!")
}

// GetRedisClient 获取Redis客户端实例
func GetRedisClient() *redis.Client {
	return RedisClient
}
