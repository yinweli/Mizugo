package contexts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestContext(t *testing.T) {
	suite.Run(t, new(SuiteContext))
}

type SuiteContext struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteContext) SetupSuite() {
	this.Change("test-contexts-context")
}

func (this *SuiteContext) TearDownSuite() {
	this.Restore()
}

func (this *SuiteContext) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteContext) TestContext() {
	assert.NotNil(this.T(), Ctx())
}
