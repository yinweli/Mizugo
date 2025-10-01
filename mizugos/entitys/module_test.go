package entitys

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestModule(t *testing.T) {
	suite.Run(t, new(SuiteModule))
}

type SuiteModule struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteModule) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-entitys-module"))
}

func (this *SuiteModule) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteModule) TestModule() {
	target := NewModule(ModuleID(1))
	this.NotNil(target)
	entity := NewEntity(EntityID(1))
	target.initialize(entity)

	this.Equal(ModuleID(1), target.ModuleID())
	this.Equal(entity, target.Entity())
}
