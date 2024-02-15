package redmos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestURI(t *testing.T) {
	suite.Run(t, new(SuiteURI))
}

type SuiteURI struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteURI) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-uri")
}

func (this *SuiteURI) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteURI) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteURI) TestRedisURIConnect() {
	_, err := RedisURI(testdata.RedisURI).Connect(ctxs.Get().Ctx())
	assert.Nil(this.T(), err)

	_, err = RedisURI("unknown://").Connect(ctxs.Get().Ctx())
	assert.NotNil(this.T(), err)

	_, err = RedisURI(testdata.RedisURIInvalid).Connect(ctxs.Get().Ctx())
	assert.NotNil(this.T(), err)
}

func (this *SuiteURI) TestRedisURIOption() {
	option, err := RedisURI("redisdb://username:password@host1:port1,host2:port2/").option()
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), option)
	assert.Equal(this.T(), "username", option.Username)
	assert.Equal(this.T(), "password", option.Password)
	assert.Equal(this.T(), []string{"host1:port1", "host2:port2"}, option.Addrs)

	option, err = RedisURI("redisdb://@host:port/").option()
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), option)
	assert.Empty(this.T(), option.Username)
	assert.Empty(this.T(), option.Password)

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
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), option)
	assert.Equal(this.T(), "name", option.ClientName)
	assert.Equal(this.T(), 1, option.DB)
	assert.Equal(this.T(), 1, option.MaxRetries)
	assert.Equal(this.T(), time.Second, option.MinRetryBackoff)
	assert.Equal(this.T(), time.Second, option.MaxRetryBackoff)
	assert.Equal(this.T(), time.Second, option.DialTimeout)
	assert.Equal(this.T(), time.Second, option.ReadTimeout)
	assert.Equal(this.T(), time.Second, option.WriteTimeout)
	assert.Equal(this.T(), false, option.ContextTimeoutEnabled)
	assert.Equal(this.T(), false, option.PoolFIFO)
	assert.Equal(this.T(), 1, option.PoolSize)
	assert.Equal(this.T(), time.Second, option.PoolTimeout)
	assert.Equal(this.T(), 1, option.MinIdleConns)
	assert.Equal(this.T(), 1, option.MaxIdleConns)
	assert.Equal(this.T(), time.Second, option.ConnMaxIdleTime)
	assert.Equal(this.T(), time.Second, option.ConnMaxLifetime)
	assert.Equal(this.T(), 1, option.MaxRedirects)
	assert.Equal(this.T(), false, option.ReadOnly)
	assert.Equal(this.T(), false, option.RouteByLatency)
	assert.Equal(this.T(), false, option.RouteRandomly)
	assert.Equal(this.T(), "master", option.MasterName)

	_, err = RedisURI("unknown://").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://@").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://@?").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://unknown@host:port/").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://username:password").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://username:password@").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://username:password@/").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://username:password@/?").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://host:port/option").option()
	assert.NotNil(this.T(), err)

	_, err = RedisURI("redisdb://username:password@host:port/?unknown").option()
	assert.NotNil(this.T(), err)
}

func (this *SuiteURI) TestRedisURIAdd() {
	target := RedisURI("redisdb://username:password@host:port/")
	assert.Equal(this.T(), target+"?"+testdata.Unknown, target.add(testdata.Unknown))
	target = "redisdb://username:password@host:port/?"
	assert.Equal(this.T(), target+"&"+testdata.Unknown, target.add(testdata.Unknown))
}

func (this *SuiteURI) TestMongoURIConnect() {
	_, err := MongoURI(testdata.MongoURI).Connect(ctxs.Get().Ctx())
	assert.Nil(this.T(), err)

	_, err = MongoURI("unknown://").Connect(ctxs.Get().Ctx())
	assert.NotNil(this.T(), err)

	_, err = MongoURI(testdata.MongoURIInvalid).Connect(ctxs.Get().Ctx())
	assert.NotNil(this.T(), err)
}

func (this *SuiteURI) TestMongoURIOption() {
	option, err := MongoURI("mongodb://username:password@sample.host:27017/?maxPoolSize=20&w=majority").option()
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), option)

	_, err = MongoURI("unknown://username:password@host1:port1,host2:port2").option()
	assert.NotNil(this.T(), err)
}
