package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEmpty(t *testing.T) {
	suite.Run(t, new(SuiteEmpty))
}

type SuiteEmpty struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEmpty) SetupSuite() {
	this.Change("test-logs-empty")
}

func (this *SuiteEmpty) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEmpty) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEmpty) TestNewEmpty() {
	assert.NotNil(this.T(), NewEmpty("", LevelDebug))
}

func (this *SuiteEmpty) TestEmpty() {
	target := NewEmpty("", LevelDebug)
	assert.Equal(this.T(), target, target.Message(""))
	assert.Equal(this.T(), target, target.KV("", 0))
	assert.Equal(this.T(), target, target.Error(nil))
	assert.Nil(this.T(), target.EndError(nil))
	target.End()
}

func (this *SuiteEmpty) TestLog() {
	target := NewEmpty("", LevelDebug)
	assert.IsType(this.T(), target, Debug(""))
	assert.IsType(this.T(), target, Info(""))
	assert.IsType(this.T(), target, Warn(""))
	assert.IsType(this.T(), target, Error(""))
	Set(NewEmpty)
	assert.IsType(this.T(), target, Debug(""))
	assert.IsType(this.T(), target, Info(""))
	assert.IsType(this.T(), target, Warn(""))
	assert.IsType(this.T(), target, Error(""))
	Set(nil)
	assert.IsType(this.T(), target, Debug(""))
	assert.IsType(this.T(), target, Info(""))
	assert.IsType(this.T(), target, Warn(""))
	assert.IsType(this.T(), target, Error(""))
}
