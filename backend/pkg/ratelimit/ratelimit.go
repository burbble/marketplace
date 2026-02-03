package ratelimit

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Max    int
	Window time.Duration
}

func Middleware(rdb *redis.Client, cfg Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ip := c.ClientIP()
		windowSec := int(cfg.Window.Seconds())
		now := time.Now().Unix()
		bucket := now / int64(windowSec)
		key := fmt.Sprintf("rl:%s:%d", ip, bucket)

		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}

		if count == 1 {
			rdb.Expire(ctx, key, cfg.Window+time.Second)
		}

		remaining := int64(cfg.Max) - count
		if remaining < 0 {
			remaining = 0
		}

		resetAt := (bucket + 1) * int64(windowSec)

		c.Header("X-RateLimit-Limit", strconv.Itoa(cfg.Max))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))

		if count > int64(cfg.Max) {
			retryAfter := resetAt - now
			c.Header("Retry-After", strconv.FormatInt(retryAfter, 10))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
