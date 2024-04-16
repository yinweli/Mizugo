package loggers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestLogmgr(t *testing.T) {
	suite.Run(t, new(SuiteLogmgr))
}

type SuiteLogmgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteLogmgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-loggers-logmgr"))
}

func (this *SuiteLogmgr) TearDownSuite() {
	trials.Restore(this.Catalog)
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
