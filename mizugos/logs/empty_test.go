package logs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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
	this.TBegin("test-logs-empty", "")
}

func (this *SuiteEmpty) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteEmpty) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteEmpty) TestEmptyLogger() {
	target := &EmptyLogger{}
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
	target.Finalize()
}

func (this *SuiteEmpty) TestEmptyStream() {
	logger := &EmptyLogger{}
	assert.Nil(this.T(), logger.Initialize())

	target := logger.Debug("")
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.Caller(0))
	assert.Equal(this.T(), target, target.KV("key", "value"))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	target.EndError(fmt.Errorf("end error"))
	target.End()

	logger.Finalize()
}
