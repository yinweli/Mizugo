package ctxs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestContext(t *testing.T) {
	suite.Run(t, new(SuiteContext))
}

type SuiteContext struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteContext) SetupSuite() {
	this.Env = testdata.EnvSetup("test-contexts-context")
}

func (this *SuiteContext) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteContext) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteContext) TestRoot() {
	Set(context.Background())
	assert.NotNil(this.T(), Get())
}

func (this *SuiteContext) TestCtx() {
	target := Ctx{}
	target.ctx, target.cancel = context.WithCancel(context.Background())
	assert.NotNil(this.T(), target.Ctx())
	withCancel := target.WithCancel()
	assert.NotNil(this.T(), withCancel.Ctx())
	withTimeout := target.WithTimeout(testdata.Timeout)
	assert.NotNil(this.T(), withTimeout.Ctx())
	withDeadline := target.WithDeadline(time.Now())
	assert.NotNil(this.T(), withDeadline.Ctx())
	target.Cancel()
	<-target.Done()
}
