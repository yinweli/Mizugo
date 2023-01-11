package contexts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestContext(t *testing.T) {
	suite.Run(t, new(SuiteContext))
}

type SuiteContext struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteContext) SetupSuite() {
	this.Change("test-contexts-context")
}

func (this *SuiteContext) TearDownSuite() {
	this.Restore()
}

func (this *SuiteContext) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteContext) TestContext() {
	assert.NotNil(this.T(), Ctx())
}
