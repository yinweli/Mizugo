package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugo/ctxs"
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
	this.Env = testdata.EnvSetup("test-redmos-major")
}

func (this *SuiteMajor) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteMajor) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(testdata.RedisURI)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Submit())
	assert.NotNil(this.T(), target.Client())
	assert.Nil(this.T(), target.SwitchDB(1))
	assert.NotNil(this.T(), target.SwitchDB(999999))
	target.DropDB()

	_, err = newMajor("")
	assert.NotNil(this.T(), err)

	ping, err := target.Client().Ping(ctxs.Get().Ctx()).Result()
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), "PONG", ping)

	target.stop()
	assert.Nil(this.T(), target.Submit())
	assert.Nil(this.T(), target.Client())
	assert.NotNil(this.T(), target.SwitchDB(1))
	target.DropDB()

	_, err = newMajor(testdata.RedisURIInvalid)
	assert.NotNil(this.T(), err)
}
