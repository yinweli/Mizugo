package entitys

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"

    "github.com/yinweli/Mizugo/testdata"
)

func TestModuler(t *testing.T) {
    suite.Run(t, new(SuiteModuler))
}

type SuiteModuler struct {
    suite.Suite
    testdata.TestEnv
}

func (this *SuiteModuler) SetupSuite() {
    this.Change("test-moduler")
}

func (this *SuiteModuler) TearDownSuite() {
    this.Restore()
}

func (this *SuiteModuler) TestNewModuler() {
    assert.NotNil(this.T(), NewModuler())
}

func (this *SuiteModuler) TestAdd() {
    target := NewModuler()

    assert.Nil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
    assert.NotNil(this.T(), target.Add(NewModule(ModuleID(1), "module")))
}

func (this *SuiteModuler) TestDel() {
    target := NewModuler()
    _ = target.Add(NewModule(ModuleID(1), "module"))

    result := target.Del(ModuleID(1))
    assert.NotNil(this.T(), result)
    assert.Equal(this.T(), ModuleID(1), result.ModuleID())
    assert.Equal(this.T(), "module", result.Name())

    result = target.Del(ModuleID(1))
    assert.Nil(this.T(), result)
}

func (this *SuiteModuler) TestGet() {
    target := NewModuler()
    _ = target.Add(NewModule(ModuleID(1), "module"))

    module := target.Get(ModuleID(1))
    assert.NotNil(this.T(), module)
    assert.Equal(this.T(), ModuleID(1), module.ModuleID())
    assert.Equal(this.T(), "module", module.Name())

    module = target.Get(ModuleID(2))
    assert.Nil(this.T(), module)
}

func (this *SuiteModuler) TestAll() {
    target := NewModuler()
    _ = target.Add(NewModule(ModuleID(1), "module1"))
    _ = target.Add(NewModule(ModuleID(2), "module2"))

    module := target.All()
    assert.Len(this.T(), module, 2)
    assert.Equal(this.T(), ModuleID(1), module[0].ModuleID())
    assert.Equal(this.T(), "module1", module[0].Name())
    assert.Equal(this.T(), ModuleID(2), module[1].ModuleID())
    assert.Equal(this.T(), "module2", module[1].Name())

}
