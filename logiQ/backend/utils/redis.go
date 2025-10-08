package utils

import (
	"backend/database"
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// SetCache 设置缓存
// key: 缓存的键
// value: 缓存的值
// expiration: 缓存的过期时间，如果为0则永不过期
func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.GetRedisClient().Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Printf("Failed to set cache for key %s: %v", key, err)
		return err
	}
	return nil
}

// GetCache 获取缓存
// key: 缓存的键
// 返回值: 缓存的值（字符串类型），如果键不存在或获取失败则返回空字符串和错误
func GetCache(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	val, err := database.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// 键不存在
			return "", nil
		}
		log.Printf("Failed to get cache for key %s: %v", key, err)
		return "", err
	}
	return val, nil
}

// DeleteCache 删除缓存
// key: 要删除的缓存键
func DeleteCache(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := database.GetRedisClient().Del(ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete cache for key %s: %v", key, err)
		return err
	}
	return nil
}