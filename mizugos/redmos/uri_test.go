package redmos

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestURI(t *testing.T) {
	suite.Run(t, new(SuiteURI))
}

type SuiteURI struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteURI) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-uri"))
}

func (this *SuiteURI) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteURI) TestRedisURIConnect() {
	_, err := RedisURI(testdata.RedisURI).Connect(context.Background())
	this.Nil(err)

	_, err = RedisURI("unknown://").Connect(context.Background())
	this.NotNil(err)

	_, err = RedisURI(testdata.RedisURIInvalid).Connect(context.Background())
	this.NotNil(err)
}

func (this *SuiteURI) TestRedisURIOption() {
	option, err := RedisURI("redisdb://username:password@host1:port1,host2:port2/").option()
	this.Nil(err)
	this.NotNil(option)
	this.Equal("username", option.Username)
	this.Equal("password", option.Password)
	this.Equal([]string{"host1:port1", "host2:port2"}, option.Addrs)

	option, err = RedisURI("redisdb://@host:port/").option()
	this.Nil(err)
	this.NotNil(option)
	this.Empty(option.Username)
	this.Empty(option.Password)

	option, err = RedisURI(
		"redisdb://username:password@host:port/?" +
			"clientName=name&" +
			"dbid=1&" +
			"maxRetries=1&" +
			"minRetryBackoff=1s&" +
			"maxRetryBackoff=1s&" +
			"dialTimeout=1s&" +
			"readTimeout=1s&" +
			"writeTimeout=1s&" +
			"contextTimeoutEnabled=false&" +
			"poolFIFO=false&" +
			"poolSize=1&" +
			"poolTimeout=1s&" +
			"minIdleConns=1&" +
			"maxIdleConns=1&" +
			"connMaxIdleTime=1s&" +
			"connMaxLifetime=1s&" +
			"maxRedirects=1&" +
			"readOnly=false&" +
			"routeByLatency=false&" +
			"routeRandomly=false&" +
			"masterName=master").option()
	this.Nil(err)
	this.NotNil(option)
	this.Equal("name", option.ClientName)
	this.Equal(1, option.DB)
	this.Equal(1, option.MaxRetries)
	this.Equal(time.Second, option.MinRetryBackoff)
	this.Equal(time.Second, option.MaxRetryBackoff)
	this.Equal(time.Second, option.DialTimeout)
	this.Equal(time.Second, option.ReadTimeout)
	this.Equal(time.Second, option.WriteTimeout)
	this.Equal(false, option.ContextTimeoutEnabled)
	this.Equal(false, option.PoolFIFO)
	this.Equal(1, option.PoolSize)
	this.Equal(time.Second, option.PoolTimeout)
	this.Equal(1, option.MinIdleConns)
	this.Equal(1, option.MaxIdleConns)
	this.Equal(time.Second, option.ConnMaxIdleTime)
	this.Equal(time.Second, option.ConnMaxLifetime)
	this.Equal(1, option.MaxRedirects)
	this.Equal(false, option.ReadOnly)
	this.Equal(false, option.RouteByLatency)
	this.Equal(false, option.RouteRandomly)
	this.Equal("master", option.MasterName)

	_, err = RedisURI("unknown://").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://@").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://@?").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://unknown@host:port/").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://username:password").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://username:password@").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://username:password@/").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://username:password@/?").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://host:port/option").option()
	this.NotNil(err)

	_, err = RedisURI("redisdb://username:password@host:port/?unknown").option()
	this.NotNil(err)
}

func (this *SuiteURI) TestRedisURIAdd() {
	target := RedisURI("redisdb://username:password@host:port/")
	this.Equal(target+"?"+testdata.Unknown, target.add(testdata.Unknown))
	target = "redisdb://username:password@host:port/?"
	this.Equal(target+"&"+testdata.Unknown, target.add(testdata.Unknown))
}

func (this *SuiteURI) TestMongoURIConnect() {
	_, err := MongoURI(testdata.MongoURI).Connect(context.Background())
	this.Nil(err)

	_, err = MongoURI("unknown://").Connect(context.Background())
	this.NotNil(err)

	_, err = MongoURI(testdata.MongoURIInvalid).Connect(context.Background())
	this.NotNil(err)
}

func (this *SuiteURI) TestMongoURIOption() {
	option, err := MongoURI("mongodb://username:password@sample.host:27017/?maxPoolSize=20&w=majority").option()
	this.Nil(err)
	this.NotNil(option)

	_, err = MongoURI("unknown://username:password@host1:port1,host2:port2").option()
	this.NotNil(err)
}
