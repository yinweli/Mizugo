package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntitymgr(t *testing.T) {
	suite.Run(t, new(SuiteEntitymgr))
}

type SuiteEntitymgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteEntitymgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-entitys-entitymgr")
}

func (this *SuiteEntitymgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteEntitymgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteEntitymgr) TestEntitymgr() {
	target := NewEntitymgr()
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Add())
	target.Clear()

	target = NewEntitymgr()
	entity := target.Add()
	assert.NotNil(this.T(), entity)
	assert.NotNil(this.T(), target.Del(entity.EntityID()))
	assert.Nil(this.T(), target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Del(entity.EntityID()))
	target.Clear()

	target = NewEntitymgr()
	assert.NotNil(this.T(), target.Add())
	assert.Equal(this.T(), 1, target.Count())
	target.Clear()
	assert.Equal(this.T(), 0, target.Count())

	target = NewEntitymgr()
	entity = target.Add()
	assert.Equal(this.T(), entity, target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Get(EntityID(2)))
	target.Clear()

	target = NewEntitymgr()
	entity1 := target.Add()
	entity2 := target.Add()
	assert.ElementsMatch(this.T(), []*Entity{entity1, entity2}, target.All())
	target.Clear()

	target = NewEntitymgr()
	assert.NotNil(this.T(), target.Add())
	assert.NotNil(this.T(), target.Add())
	assert.Equal(this.T(), 2, target.Count())
	target.Clear()
}
