package loggers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Get())
	target.Finalize()
}

func (this *SuiteEmpty) TestEmptyRetain() {
	target := &EmptyRetain{}
	assert.NotNil(this.T(), target.Clear())
	assert.NotNil(this.T(), target.Flush())
	assert.NotNil(this.T(), target.Debug(""))
	assert.NotNil(this.T(), target.Info(""))
	assert.NotNil(this.T(), target.Warn(""))
	assert.NotNil(this.T(), target.Error(""))
}

func (this *SuiteEmpty) TestEmptyStream() {
	target := &EmptyStream{retain: &EmptyRetain{}}
	assert.Equal(this.T(), target, target.Message("message"))
	assert.Equal(this.T(), target, target.KV("key", "value"))
	assert.Equal(this.T(), target, target.Caller(0))
	assert.Equal(this.T(), target, target.Error(fmt.Errorf("error")))
	assert.NotNil(this.T(), target.End())
	target.EndFlush()
}
