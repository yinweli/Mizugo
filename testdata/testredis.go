package testdata

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// TestRedis 測試redis
type TestRedis struct {
	key []string // 索引列表
}

// Key 取得索引
func (this *TestRedis) Key(key string) string {
	key = "test:" + key
	this.key = append(this.key, key)
	return key
}

// RestoreRedis 復原redis
func (this *TestRedis) RestoreRedis(ctx context.Context, client redis.UniversalClient) {
	for _, itor := range this.key {
		client.Del(ctx, itor)
	} // for
}
