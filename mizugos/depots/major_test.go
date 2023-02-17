package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
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
	target, err := newMajor(contexts.Ctx(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(contexts.Ctx(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Runner())
	assert.NotNil(this.T(), target.Client())
	target.stop()
	assert.Nil(this.T(), target.Runner())
	assert.Nil(this.T(), target.Client())

	_, err = newMajor(contexts.Ctx(), testdata.RedisURIInvalid)
	assert.NotNil(this.T(), err)
}

func (this *SuiteMajor) TestClient() {
	this.Reset()

	target, err := newMajor(contexts.Ctx(), testdata.RedisURI)
	assert.Nil(this.T(), err)
	client := target.Client()
	assert.NotNil(this.T(), client)

	key := this.Key("client")
	data := utils.RandString(testdata.RandStringLength)
	set, err := client.Set(contexts.Ctx(), key, data, testdata.RedisTimeout).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), "OK", set)

	del, err := client.Del(contexts.Ctx(), key).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), int64(1), del)

	this.RedisClear(contexts.Ctx(), client)
	target.stop()
}
