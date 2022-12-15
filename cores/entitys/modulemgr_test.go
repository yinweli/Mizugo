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
	module := newModuleTester(ModuleID(1))
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(module))
	assert.NotNil(this.T(), target.Get(module.ModuleID()))
	assert.NotNil(this.T(), target.Add(module))
}

func (this *SuiteModulemgr) TestDel() {
	module := newModuleTester(ModuleID(1))
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(module))
	assert.Equal(this.T(), module, target.Del(module.ModuleID()))
	assert.Nil(this.T(), target.Get(module.ModuleID()))
	assert.Nil(this.T(), target.Del(module.ModuleID()))
}

func (this *SuiteModulemgr) TestGet() {
	module := newModuleTester(ModuleID(1))
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(module))
	assert.Equal(this.T(), module, target.Get(module.ModuleID()))
	assert.Nil(this.T(), target.Get(ModuleID(2)))
}

func (this *SuiteModulemgr) TestAll() {
	module1 := newModuleTester(ModuleID(1))
	module2 := newModuleTester(ModuleID(2))
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(module1))
	assert.Nil(this.T(), target.Add(module2))
	assert.ElementsMatch(this.T(), []Moduler{module1, module2}, target.All())
}

func (this *SuiteModulemgr) TestCount() {
	target := NewModulemgr()

	assert.Nil(this.T(), target.Add(newModuleTester(ModuleID(1))))
	assert.Nil(this.T(), target.Add(newModuleTester(ModuleID(2))))
	assert.Equal(this.T(), 2, target.Count())
}