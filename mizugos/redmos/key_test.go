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
	testdata.Env
}

func (this *SuiteFormat) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-redmos-format")
}

func (this *SuiteFormat) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteFormat) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteFormat) TestMajorKey() {
	assert.Equal(this.T(), "A:B:C", MajorKey("A", "B", "C"))
}
