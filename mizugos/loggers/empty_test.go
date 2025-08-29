package loggers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestEmpty(t *testing.T) {
	suite.Run(t, new(SuiteEmpty))
}

type SuiteEmpty struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteEmpty) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-loggers-empty"))
}

func (this *SuiteEmpty) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteEmpty) TestEmptyLogger() {
	target := &EmptyLogger{}
	this.Nil(target.Initialize())
	target.fail = true
	this.NotNil(target.Initialize())
	this.NotNil(target.Get())
	target.Finalize()
}

func (this *SuiteEmpty) TestEmptyRetain() {
	target := &EmptyRetain{}
	this.NotNil(target.Clear())
	this.NotNil(target.Flush())
	this.NotNil(target.Debug(""))
	this.NotNil(target.Info(""))
	this.NotNil(target.Warn(""))
	this.NotNil(target.Error(""))
}

func (this *SuiteEmpty) TestEmptyStream() {
	target := &EmptyStream{}
	this.Equal(target, target.Message("message"))
	this.Equal(target, target.KV("key", "value"))
	this.Equal(target, target.Caller(0))
	this.Equal(target, target.Error(fmt.Errorf("error")))
	this.NotNil(target.End())
	target.EndFlush()
}
