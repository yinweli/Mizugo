package depots

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
	testdata.TestLeak
}

func (this *SuiteFormat) SetupSuite() {
	this.Change("test-depots-format")
}

func (this *SuiteFormat) TearDownSuite() {
	this.Restore()
}

func (this *SuiteFormat) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteFormat) TestFormatField() {
	assert.Equal(this.T(), "abc", FormatField("ABC"))
}

func (this *SuiteFormat) TestFormatKey() {
	assert.Equal(this.T(), "a:b:c", FormatKey("A", "B", "C"))
}
