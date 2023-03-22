package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestFormat(t *testing.T) {
	suite.Run(t, new(SuiteFormat))
}

type SuiteFormat struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteFormat) SetupSuite() {
	this.TBegin("test-redmos-format", "")
}

func (this *SuiteFormat) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteFormat) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteFormat) TestFormatField() {
	assert.Equal(this.T(), "abc", FormatField("ABC"))
}

func (this *SuiteFormat) TestFormatKey() {
	assert.Equal(this.T(), "a:b:c", FormatKey("A", "B", "C"))
}
