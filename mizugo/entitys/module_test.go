package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestModule(t *testing.T) {
	suite.Run(t, new(SuiteModule))
}

type SuiteModule struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteModule) SetupSuite() {
	this.Env = testdata.EnvSetup("test-entitys-module")
}

func (this *SuiteModule) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteModule) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteModule) TestModule() {
	target := NewModule(ModuleID(1))
	assert.NotNil(this.T(), target)
	entity := NewEntity(EntityID(1))
	target.initialize(entity)

	assert.Equal(this.T(), ModuleID(1), target.ModuleID())
	assert.Equal(this.T(), entity, target.Entity())
}
