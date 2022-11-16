package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestModule(t *testing.T) {
	suite.Run(t, new(SuiteModule))
}

type SuiteModule struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteModule) SetupSuite() {
	this.Change("test-module")
}

func (this *SuiteModule) TearDownSuite() {
	this.Restore()
}

func (this *SuiteModule) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteModule) TestNewModule() {
	assert.NotNil(this.T(), NewModule(0, ""))
}

func (this *SuiteModule) TestModule() {
	target := NewModule(ModuleID(1), "module")
	target.Host(NewEntity(EntityID(1), "entity"))

	assert.Equal(this.T(), ModuleID(1), target.ModuleID())
	assert.Equal(this.T(), "module", target.Name())
	assert.Equal(this.T(), EntityID(1), target.Entity().EntityID())
	assert.Equal(this.T(), "entity", target.Entity().Name())
}
