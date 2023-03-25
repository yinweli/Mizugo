package logs

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
	testdata.EnvSetup(&this.Env, "test-logs-logmgr")
}

func (this *SuiteLogmgr) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteLogmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteLogmgr) TestNewLogmgr() {
	assert.NotNil(this.T(), NewLogmgr())
}

func (this *SuiteLogmgr) TestLogmgr() {
	target := NewLogmgr()
	name := "log"
	assert.Nil(this.T(), target.Add(name, newLoggerTester(true)))
	assert.NotNil(this.T(), target.Get(name))
	assert.Nil(this.T(), target.Get(testdata.Unknown))
	target.Finalize()

	assert.NotNil(this.T(), target.Add("", newLoggerTester(true)))
	assert.Nil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log2", nil))
	assert.NotNil(this.T(), target.Add("log3", newLoggerTester(false)))
}
