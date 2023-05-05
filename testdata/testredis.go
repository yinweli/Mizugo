package testdata

import (
	"context"
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/redis/go-redis/v9"
)

// RedisCompare 在redis中比對資料是否相同
func RedisCompare[T any](client redis.Cmdable, key string, expected *T, cmpOpt ...cmp.Option) bool {
	result, err := client.Get(context.Background(), key).Result()

	if err != nil {
		return false
	} // if

	actual := new(T)

	if json.Unmarshal([]byte(result), actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, cmpOpt...)
}
