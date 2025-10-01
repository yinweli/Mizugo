package trials

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestRedis(t *testing.T) {
	suite.Run(t, new(SuiteRedis))
}

type SuiteRedis struct {
	suite.Suite
}

func (this *SuiteRedis) TestRedisExist() {
	client := newRedis()
	client.Set(context.Background(), "redis exist", "1", 0)
	this.True(RedisExist(client, "redis exist"))
	this.False(RedisExist(client, testdata.Unknown))
	client.Del(context.Background(), "redis exist")
}

func (this *SuiteRedis) TestRedisEqual() {
	client := newRedis()
	client.Set(context.Background(), "redis compare", "{ \"Value\": 1 }", 0)
	client.Set(context.Background(), "redis compare?", "1234567890123456789012345678901234567890", 0)
	this.True(RedisEqual[testRedis](client, "redis compare", &testRedis{Value: 1}))
	this.False(RedisEqual[testRedis](client, testdata.Unknown, &testRedis{}))
	this.False(RedisEqual[testRedis](client, "redis compare?", &testRedis{}))
	this.False(RedisEqual[testRedis](client, "redis compare", &testRedis{Value: 2}))
	client.Del(context.Background(), "redis compare")
	client.Del(context.Background(), "redis compare?")
}

func (this *SuiteRedis) TestRedisListEqual() {
	client := newRedis()
	client.RPush(context.Background(), "redis compare list", "{ \"Value\": 1 }", "{ \"Value\": 2 }")
	client.RPush(context.Background(), "redis compare list?", "1234567890123456789012345678901234567890")
	client.Set(context.Background(), "redis compare list@", testdata.Unknown, 0)
	this.True(RedisListEqual[testRedis](client, "redis compare list", []*testRedis{
		{Value: 1},
		{Value: 2},
	}))
	this.False(RedisListEqual[testRedis](client, "redis compare list@", []*testRedis{}))
	this.False(RedisListEqual[testRedis](client, "redis compare list?", []*testRedis{}))
	this.False(RedisListEqual[testRedis](client, "redis compare list", []*testRedis{
		{Value: 2},
		{Value: 1},
	}))
	client.Del(context.Background(), "redis compare list")
	client.Del(context.Background(), "redis compare list?")
	client.Del(context.Background(), "redis compare list@")
}

func newRedis() redis.Cmdable {
	option := &redis.UniversalOptions{}
	option.Addrs = append(option.Addrs, testdata.RedisIP)
	client := redis.NewUniversalClient(option)
	return client
}

type testRedis struct {
	Value int
}
