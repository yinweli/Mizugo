package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntity(t *testing.T) {
	suite.Run(t, new(SuiteEntity))
}

type SuiteEntity struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntity) SetupSuite() {
	this.Change("test-entity")
}

func (this *SuiteEntity) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntity) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntity) TestNewEntity() {
	assert.NotNil(this.T(), NewEntity(0, ""))
}

func (this *SuiteEntity) TestEntity() {
	target := NewEntity(EntityID(1), "entity")

	assert.Equal(this.T(), EntityID(1), target.EntityID())
	assert.Equal(this.T(), "entity", target.Name())
}

func (this *SuiteEntity) TestAddModule() {
	target := NewEntity(EntityID(1), "entity")

	assert.Nil(this.T(), target.AddModule(NewModule(ModuleID(1), "module")))
}

func (this *SuiteEntity) TestDelModule() {
	target := NewEntity(EntityID(1), "entity")
	_ = target.AddModule(NewModule(ModuleID(1), "module"))

	assert.NotNil(this.T(), target.DelModule(ModuleID(1)))
	assert.Nil(this.T(), target.DelModule(ModuleID(1)))
}

func (this *SuiteEntity) TestGetModule() {
	target := NewEntity(EntityID(1), "entity")
	_ = target.AddModule(NewModule(ModuleID(1), "module"))

	assert.NotNil(this.T(), target.GetModule(ModuleID(1)))
	assert.Nil(this.T(), target.GetModule(ModuleID(2)))
}
