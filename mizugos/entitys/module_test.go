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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteModule) SetupSuite() {
	this.Change("test-entitys-module")
}

func (this *SuiteModule) TearDownSuite() {
	this.Restore()
}

func (this *SuiteModule) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteModule) TestNewModule() {
	assert.NotNil(this.T(), NewModule(ModuleID(0)))
}

func (this *SuiteModule) TestModule() {
	target := NewModule(ModuleID(1))
	entity := NewEntity(EntityID(1))
	target.initialize(entity)

	assert.Equal(this.T(), ModuleID(1), target.ModuleID())
	assert.Equal(this.T(), entity, target.Entity())
}
