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

func (this *SuiteLogmgr) TestInitialize() {
	target := NewLogmgr()
	assert.Nil(this.T(), target.Initialize(newLoggerTester(true)))
	assert.NotNil(this.T(), target.Initialize(newLoggerTester(false)))
	target.Finalize()
}

func (this *SuiteLogmgr) TestLog() {
	target := NewLogmgr()
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
	assert.Nil(this.T(), target.Initialize(newLoggerTester(true)))
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
	target.Finalize()
}
