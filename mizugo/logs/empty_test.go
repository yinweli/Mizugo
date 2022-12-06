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
	Initialize(&EmptyLogger{})
	assert.NotNil(this.T(), Debug(""))
	assert.NotNil(this.T(), Info(""))
	assert.NotNil(this.T(), Warn(""))
	assert.NotNil(this.T(), Error(""))

	Initialize(nil)
	assert.NotNil(this.T(), Debug(""))
	assert.NotNil(this.T(), Info(""))
	assert.NotNil(this.T(), Warn(""))
	assert.NotNil(this.T(), Error(""))
	Finalize()
}

func (this *SuiteEmpty) TestEmptyStream() {
	Initialize(&EmptyLogger{})

	target := Debug("log")
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.KV("key", "value"))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	assert.NotNil(this.T(), target.EndError(fmt.Errorf("end error")))
	target.End()

	Finalize()
	Initialize(nil)
}
