package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestModulemgr(t *testing.T) {
	suite.Run(t, new(SuiteModulemgr))
}

type SuiteModulemgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteModulemgr) SetupSuite() {
	this.Change("test-entitys-modulemgr")
}

func (this *SuiteModulemgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteModulemgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteModulemgr) TestNewModulemgr() {
	assert.NotNil(this.T(), NewModulemgr())
}

func (this *SuiteModulemgr) TestAdd() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
	assert.NotNil(this.T(), target.Get(ModuleID(1)))

	assert.NotNil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
}

func (this *SuiteModulemgr) TestDel() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
	module := target.Del(ModuleID(1))
	assert.NotNil(this.T(), module)
	assert.Equal(this.T(), ModuleID(1), module.ModuleID())
	assert.Equal(this.T(), "module", module.Name())
	assert.Nil(this.T(), target.Get(ModuleID(1)))

	assert.Nil(this.T(), target.Del(ModuleID(1)))
}

func (this *SuiteModulemgr) TestGet() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
	module := target.Get(ModuleID(1))
	assert.NotNil(this.T(), module)
	assert.Equal(this.T(), ModuleID(1), module.ModuleID())
	assert.Equal(this.T(), "module", module.Name())

	assert.Nil(this.T(), target.Get(ModuleID(2)))
}

func (this *SuiteModulemgr) TestAll() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module1")))
	assert.Nil(this.T(), target.Add(NewModule(ModuleID(2), "module2")))
	module := target.All()
	assert.Len(this.T(), module, 2)
	assert.Equal(this.T(), ModuleID(1), module[0].ModuleID())
	assert.Equal(this.T(), "module1", module[0].Name())
	assert.Equal(this.T(), ModuleID(2), module[1].ModuleID())
	assert.Equal(this.T(), "module2", module[1].Name())
}

func (this *SuiteModulemgr) TestCount() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module1")))
	assert.Nil(this.T(), target.Add(NewModule(ModuleID(2), "module2")))
	assert.Equal(this.T(), 2, target.Count())
}
