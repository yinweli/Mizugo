package entitys

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestModulemap(t *testing.T) {
	suite.Run(t, new(SuiteModulemap))
}

type SuiteModulemap struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteModulemap) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-entitys-modulemap"))
}

func (this *SuiteModulemap) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteModulemap) TestModulemap() {
	module1 := newTestModule(true, true, ModuleID(1))
	module2 := newTestModule(true, true, ModuleID(2))
	target := NewModulemap()
	this.NotNil(target)
	this.Nil(target.Add(module1))
	this.NotNil(target.Get(module1.ModuleID()))
	this.NotNil(target.Add(module1))

	target = NewModulemap()
	this.Nil(target.Add(module1))
	this.Equal(module1, target.Del(module1.ModuleID()))
	this.Nil(target.Get(module1.ModuleID()))
	this.Nil(target.Del(module1.ModuleID()))

	target = NewModulemap()
	this.Nil(target.Add(module1))
	this.Equal(module1, target.Get(module1.ModuleID()))
	this.Nil(target.Get(ModuleID(2)))

	target = NewModulemap()
	this.Nil(target.Add(module1))
	this.Nil(target.Add(module2))
	this.ElementsMatch([]Moduler{module1, module2}, target.All())
	this.Equal(2, target.Count())
}
