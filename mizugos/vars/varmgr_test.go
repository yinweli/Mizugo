package vars

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestVarmgr(t *testing.T) {
	suite.Run(t, new(SuiteVarmgr))
}

type SuiteVarmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteVarmgr) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-vars-varmgr")
}

func (this *SuiteVarmgr) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteVarmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteVarmgr) TestNewVarmgr() {
	assert.NotNil(this.T(), NewVarmgr())
}

func (this *SuiteVarmgr) TestVarmgr() {
	target := NewVarmgr()
	name := "var"
	data := "data"

	assert.Nil(this.T(), target.Get(name))
	target.Set(name, data)
	target.Reset()
	assert.Nil(this.T(), target.Get(name))
	target.Set(name, data)
	assert.Equal(this.T(), data, target.Get(name))
	target.Del(name)
	assert.Nil(this.T(), target.Get(name))
}
