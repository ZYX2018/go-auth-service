package resources

import (
	"github.com/go-redis/redis/v8"
	"go-auth-service/config"
)

func InitRedisResource(appConfig *config.AppConfig) *redis.Client {
	// 配置Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     appConfig.Redis.Addr,
		Password: appConfig.Redis.Password,
		DB:       appConfig.Redis.DB,
	})
	return rdb
}
