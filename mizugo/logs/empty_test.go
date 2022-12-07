package logs

import (
	"fmt"
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

func (this *SuiteEmpty) TestEmptyLogger() {
	target := &EmptyLogger{}
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.New("", LevelDebug))
	target.Finalize()
}

func (this *SuiteEmpty) TestEmptyStream() {
	logger := &EmptyLogger{}
	assert.Nil(this.T(), logger.Initialize())

	target := logger.New("log", LevelDebug)
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.KV("key", "value"))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	assert.NotNil(this.T(), target.EndError(fmt.Errorf("end error")))
	target.End()

	logger.Finalize()
}
