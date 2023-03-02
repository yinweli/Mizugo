package depots

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
	testdata.TestEnv
	testdata.TestLeak
	testdata.TestDB
}

func (this *SuiteMajor) SetupSuite() {
	this.Change("test-depots-major")
}

func (this *SuiteMajor) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMajor) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMajor) TestNewMajor() {
	target, err := newMajor(ctxs.Root(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(ctxs.Root(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())

	_, err = newMajor(ctxs.Root(), testdata.RedisURIInvalid)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMajor) TestClient() {
	target, err := newMajor(ctxs.Root(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	client := target.Client()
	assert.NotNil(this.T(), client)

	key := this.Key("major client")
	data := utils.RandString(testdata.RandStringLength)
	set, err := client.Set(ctxs.RootCtx(), key, data, testdata.RedisTimeout).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), RedisOk, set)

	del, err := client.Del(ctxs.RootCtx(), key).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), int64(1), del)

	this.RedisClear(ctxs.RootCtx(), client)
	target.stop()
}

func BenchmarkMajorSet(b *testing.B) {
	target, _ := newMajor(ctxs.Root(), testdata.RedisURI)
	submit := target.Submit()

	for i := 0; i < b.N; i++ {
		value := utils.RandString(testdata.RandStringLength)
		_, _ = submit.Set(ctxs.RootCtx(), value, value, testdata.RedisTimeout).Result()
	} // for

	_, _ = submit.Exec(ctxs.RootCtx())
}

func BenchmarkMajorGet(b *testing.B) {
	target, _ := newMajor(ctxs.Root(), testdata.RedisURI)
	submit := target.Submit()
	value := utils.RandString(testdata.RandStringLength)
	_, _ = submit.Set(ctxs.RootCtx(), value, value, 0).Result()

	for i := 0; i < b.N; i++ {
		_, _ = submit.Get(ctxs.RootCtx(), value).Result()
	} // for

	_, _ = submit.Del(ctxs.RootCtx(), value).Result()
	_, _ = submit.Exec(ctxs.RootCtx())
}
