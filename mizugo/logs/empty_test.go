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
	this.Change("test-empty-log")
}

func (this *SuiteEmpty) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEmpty) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEmpty) TestLog() {
	assert.IsType(this.T(), &Empty{}, Debug(""))
	assert.IsType(this.T(), &Empty{}, Info(""))
	assert.IsType(this.T(), &Empty{}, Warn(""))
	assert.IsType(this.T(), &Empty{}, Error(""))
	Set(NewEmpty)
	assert.IsType(this.T(), &Empty{}, Debug(""))
	assert.IsType(this.T(), &Empty{}, Info(""))
	assert.IsType(this.T(), &Empty{}, Warn(""))
	assert.IsType(this.T(), &Empty{}, Error(""))
	Set(nil)
	assert.IsType(this.T(), &Empty{}, Debug(""))
	assert.IsType(this.T(), &Empty{}, Info(""))
	assert.IsType(this.T(), &Empty{}, Warn(""))
	assert.IsType(this.T(), &Empty{}, Error(""))
}

func (this *SuiteEmpty) TestEmpty() {
	target := &Empty{}
	assert.Equal(this.T(), target, target.Message(""))
	assert.Equal(this.T(), target, target.KV("", 0))
	assert.Equal(this.T(), target, target.Error(nil))
	assert.Nil(this.T(), target.EndError(nil))
	target.End()
}
