package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMajor(t *testing.T) {
	suite.Run(t, new(SuiteMajor))
}

type SuiteMajor struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
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
	target, err := newMajor(contexts.Ctx(), "redisdb://127.0.0.1:6379/")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
}

func (this *SuiteMajor) TestMajor() {
	target, err := newMajor(contexts.Ctx(), "redisdb://127.0.0.1:6379/")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Runner())
	assert.NotNil(this.T(), target.Client())
	target.stop()
	assert.Nil(this.T(), target.Runner())
	assert.Nil(this.T(), target.Client())

	_, err = newMajor(contexts.Ctx(), "redisdb://127.0.0.1:10001/?dialTimeout=1s")
	assert.NotNil(this.T(), err)
}
