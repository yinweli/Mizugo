package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestLogger(t *testing.T) {
	suite.Run(t, new(SuiteLogger))
}

type SuiteLogger struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteLogger) SetupSuite() {
	this.Change("test-logger")
}

func (this *SuiteLogger) TearDownSuite() {
	this.Restore()
}

func (this *SuiteLogger) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteLogger) TestLogger() {
	New = func() Logger {
		return nil
	}
	assert.Nil(this.T(), Begin())
	New = nil
	assert.NotNil(this.T(), Begin())
}

func (this *SuiteLogger) TestEmpty() {
	target := &empty{}
	assert.Equal(this.T(), target, target.Message(""))
	assert.Equal(this.T(), target, target.KeyValue("", 0))
	assert.Equal(this.T(), target, target.Error(nil))
	assert.Nil(this.T(), target.EndError(nil))
}
