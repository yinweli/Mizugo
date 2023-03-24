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
	testdata.EnvSetup(&this.Env, "test-contexts-context")
}

func (this *SuiteContext) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteContext) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteContext) TestRoot() {
	assert.NotNil(this.T(), Root())
	assert.NotNil(this.T(), RootCtx())
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
