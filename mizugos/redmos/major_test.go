package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMajor(t *testing.T) {
	suite.Run(t, new(SuiteMajor))
}

type SuiteMajor struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteMajor) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-major")
}

func (this *SuiteMajor) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteMajor) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMajor) TestNewMajor() {
	target, err := newMajor(ctxs.RootCtx(), testdata.RedisURI, true)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)

	_, err = newMajor(ctxs.RootCtx(), "", true)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(ctxs.RootCtx(), testdata.RedisURI, true)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	ping, err := target.Client().Ping(ctxs.RootCtx()).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), "PONG", ping)

	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())

	_, err = newMajor(ctxs.RootCtx(), testdata.RedisURIInvalid, true)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMajor) TestUsedKey() {
	target, err := newMajor(ctxs.RootCtx(), testdata.RedisURI, true)
	assert.Nil(this.T(), err)
	client := target.Client()
	assert.NotNil(this.T(), client)

	data := utils.RandString(testdata.RandStringLength)
	_, _ = client.Set(ctxs.RootCtx(), "index1", data, testdata.RedisTimeout).Result()
	_, _ = client.Set(ctxs.RootCtx(), "index2", data, testdata.RedisTimeout).Result()
	_, _ = client.Set(ctxs.RootCtx(), "index3", data, testdata.RedisTimeout).Result()
	assert.Equal(this.T(), []string{"index1", "index2", "index3"}, target.UsedKey())
	target.stop()
}

func BenchmarkMajorSet(b *testing.B) {
	target, _ := newMajor(ctxs.RootCtx(), testdata.RedisURI, false)
	submit := target.Submit()

	for i := 0; i < b.N; i++ {
		value := utils.RandString(testdata.RandStringLength)
		_, _ = submit.Set(ctxs.RootCtx(), value, value, testdata.RedisTimeout).Result()
	} // for

	_, _ = submit.Exec(ctxs.RootCtx())
}

func BenchmarkMajorGet(b *testing.B) {
	target, _ := newMajor(ctxs.RootCtx(), testdata.RedisURI, false)
	submit := target.Submit()
	value := utils.RandString(testdata.RandStringLength)
	_, _ = submit.Set(ctxs.RootCtx(), value, value, 0).Result()

	for i := 0; i < b.N; i++ {
		_, _ = submit.Get(ctxs.RootCtx(), value).Result()
	} // for

	_, _ = submit.Del(ctxs.RootCtx(), value).Result()
	_, _ = submit.Exec(ctxs.RootCtx())
}
