package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	module1 := newModuleTester(true, true, ModuleID(1))
	module2 := newModuleTester(true, true, ModuleID(2))
	target := NewModulemap()
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Add(module1))
	assert.NotNil(this.T(), target.Get(module1.ModuleID()))
	assert.NotNil(this.T(), target.Add(module1))

	target = NewModulemap()
	assert.Nil(this.T(), target.Add(module1))
	assert.Equal(this.T(), module1, target.Del(module1.ModuleID()))
	assert.Nil(this.T(), target.Get(module1.ModuleID()))
	assert.Nil(this.T(), target.Del(module1.ModuleID()))

	target = NewModulemap()
	assert.Nil(this.T(), target.Add(module1))
	assert.Equal(this.T(), module1, target.Get(module1.ModuleID()))
	assert.Nil(this.T(), target.Get(ModuleID(2)))

	target = NewModulemap()
	assert.Nil(this.T(), target.Add(module1))
	assert.Nil(this.T(), target.Add(module2))
	assert.ElementsMatch(this.T(), []Moduler{module1, module2}, target.All())
	assert.Equal(this.T(), 2, target.Count())
}
