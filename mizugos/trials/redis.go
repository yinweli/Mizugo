package trials

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/redis/go-redis/v9"
)

// RedisExist 檢查索引是否存在
func RedisExist(client redis.Cmdable, key string) bool {
	result, err := client.Exists(context.Background(), key).Result()

	if err != nil {
		fmt.Printf("redis not exist: key %v: %v\n", key, err)
		return false
	} // if

	if result == 0 {
		fmt.Printf("redis not exist: key %v: not exist\n", key)
		return false
	} // if

	return true
}

// RedisEqual 比對資料是否符合預期
func RedisEqual[T any](client redis.Cmdable, key string, expected *T, option ...cmp.Option) bool {
	result, err := client.Get(context.Background(), key).Result()

	if err != nil {
		fmt.Printf("redis not equal: key %v: %v\n", key, err)
		return false
	} // if

	actual := new(T)

	if err = json.Unmarshal([]byte(result), actual); err != nil {
		raw := result

		if len(raw) > 32 { //nolint:mnd
			raw = raw[:32] + "..."
		} // if

		fmt.Printf("redis not equal: key %v: raw: %v: %v\n", key, raw, err)
		return false
	} // if

	if cmp.Equal(expected, actual, option...) == false {
		fmt.Printf("redis not equal: key %v:\n", key)
		fmt.Println("  expected:")
		fmt.Printf("    %+v\n", expected)
		fmt.Println("  actual:")
		fmt.Printf("    %+v\n", actual)
		return false
	} // if

	return true
}

// RedisListEqual 比對資料列表是否符合預期
func RedisListEqual[T any](client redis.Cmdable, key string, expected []*T, option ...cmp.Option) bool {
	result, err := client.LRange(context.Background(), key, 0, -1).Result()

	if err != nil {
		fmt.Printf("redis not equal: key %v: %v\n", key, err)
		return false
	} // if

	actual := []*T{}

	for _, itor := range result {
		a := new(T)

		if err = json.Unmarshal([]byte(itor), a); err != nil {
			raw := itor

			if len(raw) > 32 { //nolint:mnd
				raw = raw[:32] + "..."
			} // if

			fmt.Printf("redis not equal: key %v: raw: %v: %v\n", key, raw, err)
			return false
		} // if

		actual = append(actual, a)
	} // for

	if cmp.Equal(expected, actual, option...) == false {
		fmt.Printf("redis not equal: key %v\n", key)
		fmt.Println("  expected:")

		for _, itor := range expected {
			fmt.Printf("    %+v\n", itor)
		} // for

		fmt.Println("  actual:")

		for _, itor := range actual {
			fmt.Printf("    %+v\n", itor)
		} // for

		return false
	} // if

	return true
}
