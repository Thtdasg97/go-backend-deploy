package cache

import (
	"context"
	"fmt"
	"go-backend/internal/common/env"
	"go-backend/internal/common/helpers"
	"log"
	"time"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	memoryStore *persist.MemoryStore
	redisStore  *persist.RedisStore
}

func NewCache(env *env.Env) *Cache {
	memoryStore := persist.NewMemoryStore(1 * time.Minute)

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network:  "tcp",
		Addr:     env.RedisAddr,
		Password: env.RedisPass,
	}))

	// check xem kết nối thành công ko
	statusCmd := redisStore.RedisClient.Ping(context.Background())
	err := statusCmd.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[REDIS] Connection To Redis Successfully", statusCmd.String())

	return &Cache{
		memoryStore: memoryStore,
		redisStore:  redisStore,
	}
}

// cache by memory

func (c *Cache) CacheMemory(defaultExpire time.Duration) gin.HandlerFunc {
	return cache.CacheByRequestURI(c.memoryStore,
		defaultExpire,
		cache.WithCacheStrategyByRequest(customKey))
}

// cache by redis
func (c *Cache) CacheRedis(defaultExpire time.Duration) gin.HandlerFunc {
	return cache.CacheByRequestURI(c.redisStore,
		defaultExpire,
		cache.WithCacheStrategyByRequest(customKey))
}

// này là đang implement cái gì vậy?
func customKey(ctx *gin.Context) (bool, cache.Strategy) {
	user, err := helpers.GetUser(ctx)

	key := fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.Request.URL.RequestURI())

	if err == nil && user != nil {
		key = fmt.Sprintf("%d:%s:%s",
			user.ID, ctx.Request.Method, ctx.Request.URL.RequestURI())
	}

	return true, cache.Strategy{
		CacheKey: key,
	}
}
