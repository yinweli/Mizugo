package testdata

import (
	"context"
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/redis/go-redis/v9"
)

// RedisClear 清除redis
func RedisClear(ctx context.Context, client redis.Cmdable, key []string) {
	for _, itor := range key {
		client.Del(ctx, itor)
	} // for
}

// RedisCompare 在redis中比對資料是否相同
func RedisCompare[T any](ctx context.Context, client redis.Cmdable, key string, expected *T, cmpOpt ...cmp.Option) bool {
	result, err := client.Get(ctx, key).Result()

	if err != nil {
		return false
	} // if

	actual := new(T)

	if json.Unmarshal([]byte(result), actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, cmpOpt...)
}
