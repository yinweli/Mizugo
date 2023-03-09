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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteLogmgr) SetupSuite() {
	this.Change("test-logs-logmgr")
}

func (this *SuiteLogmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteLogmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteLogmgr) TestNewLogmgr() {
	assert.NotNil(this.T(), NewLogmgr())
}

func (this *SuiteLogmgr) TestLogmgr() {
	target := NewLogmgr()
	name := "log"
	assert.Nil(this.T(), target.Add(name, newLoggerTester(true)))
	assert.NotNil(this.T(), target.Debug(name, name))
	assert.NotNil(this.T(), target.Info(name, name))
	assert.NotNil(this.T(), target.Warn(name, name))
	assert.NotNil(this.T(), target.Error(name, name))
	target.Finalize()

	assert.NotNil(this.T(), target.Add("", newLoggerTester(true)))
	assert.Nil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log1", newLoggerTester(true)))
	assert.NotNil(this.T(), target.Add("log2", nil))
	assert.NotNil(this.T(), target.Add("log3", newLoggerTester(false)))
}
