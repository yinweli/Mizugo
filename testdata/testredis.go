package testdata

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisClear 清除redis
func RedisClear(ctx context.Context, client redis.Cmdable, key []string) {
	for _, itor := range key {
		client.Del(ctx, itor)
	} // for
}
