package loggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestLogmgr(t *testing.T) {
	suite.Run(t, new(SuiteLogmgr))
}

type SuiteLogmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteLogmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-loggers-logmgr")
}

func (this *SuiteLogmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteLogmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteLogmgr) TestLogmgr() {
	target := NewLogmgr()
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Add("log", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Get("log"))
	assert.Nil(this.T(), target.Get(testdata.Unknown))
	target.Finalize()

	assert.NotNil(this.T(), target.Add("", newLoggerTester(true)))
	assert.Nil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log2", nil))
	assert.NotNil(this.T(), target.Add("log3", newLoggerTester(false)))
}
