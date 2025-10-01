package loggers

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.NotNil(target)
	target.Finalize()
}

func (this *SuiteLogmgr) TestAdd() {
	target := NewLogmgr()
	this.Nil(target.Add("log", &EmptyLogger{}))
	this.NotNil(target.Get("log"))
	this.Nil(target.Get(testdata.Unknown))
	this.NotNil(target.Add("", &EmptyLogger{}))
	this.NotNil(target.Add("log", &EmptyLogger{}))
	this.NotNil(target.Add("log nil", nil))
	this.NotNil(target.Add("log fail", &EmptyLogger{fail: true}))
	target.Finalize()
}
