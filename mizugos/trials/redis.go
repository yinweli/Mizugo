package trials

import (
	"context"
	"encoding/json"

	"github.com/google/go-cmp/cmp"
	"github.com/redis/go-redis/v9"
)

// RedisExist 在redis中資料是否存在
func RedisExist(client redis.Cmdable, key string) bool {
	result, err := client.Exists(context.Background(), key).Result()

	if err != nil {
		return false
	} // if

	return result > 0
}

// RedisCompare 在redis中比對資料是否相同
func RedisCompare[T any](client redis.Cmdable, key string, expected *T, option ...cmp.Option) bool {
	result, err := client.Get(context.Background(), key).Result()

	if err != nil {
		return false
	} // if

	actual := new(T)

	if json.Unmarshal([]byte(result), actual) != nil {
		return false
	} // if

	return cmp.Equal(expected, actual, option...)
}
