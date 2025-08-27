package trials

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
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
	key := "exist"
	client := newRedis()
	client.Set(context.Background(), key, "1", 0)
	assert.True(this.T(), RedisExist(client, key))
	assert.False(this.T(), RedisExist(client, testdata.Unknown))
	client.Del(context.Background(), key)
}

func (this *SuiteRedis) TestRedisCompare() {
	key := "compare"
	client := newRedis()
	client.Set(context.Background(), key, "{ \"Value\": 1 }", 0)
	assert.True(this.T(), RedisCompare[testRedis](client, key, &testRedis{Value: 1}))
	assert.False(this.T(), RedisCompare[testRedis](client, testdata.Unknown, nil))
	client.Del(context.Background(), key)
}

func (this *SuiteRedis) TestRedisCompareList() {
	key := "compareList"
	client := newRedis()
	client.RPush(context.Background(), key, "{ \"Value\": 1 }", "{ \"Value\": 2 }")
	assert.True(this.T(), RedisCompareList[testRedis](client, key, []*testRedis{
		{Value: 1},
		{Value: 2},
	}))
	assert.False(this.T(), RedisCompareList[testRedis](client, testdata.Unknown, nil))
	client.Del(context.Background(), key)
}

func newRedis() redis.Cmdable {
	option := &redis.UniversalOptions{}
	option.Addrs = append(option.Addrs, testdata.RedisIP)
	client := redis.NewUniversalClient(option)
	return client
}

type testRedis struct {
	Value int `json:"value"`
}
