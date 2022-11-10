package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntityr(t *testing.T) {
	suite.Run(t, new(SuiteEntityr))
}

type SuiteEntityr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntityr) SetupSuite() {
	this.Change("test-entityr")
}

func (this *SuiteEntityr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntityr) TestNewEntityr() {
	assert.NotNil(this.T(), NewEntityr())
}

func (this *SuiteEntityr) TestAdd() {
	target := NewEntityr()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.NotNil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
}

func (this *SuiteEntityr) TestDel() {
	target := NewEntityr()
	_ = target.Add(NewEntity(EntityID(1), "entity"))

	entity := target.Del(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	entity = target.Del(EntityID(2))
	assert.Nil(this.T(), entity)
}

func (this *SuiteEntityr) TestGet() {
	target := NewEntityr()
	_ = target.Add(NewEntity(EntityID(1), "entity"))

	entity := target.Get(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	entity = target.Get(EntityID(2))
	assert.Nil(this.T(), entity)
}

func (this *SuiteEntityr) TestAll() {
	target := NewEntityr()
	_ = target.Add(NewEntity(EntityID(1), "entity1"))
	_ = target.Add(NewEntity(EntityID(2), "entity2"))

	entity := target.All()
	assert.Len(this.T(), entity, 2)
	assert.Equal(this.T(), EntityID(1), entity[0].EntityID())
	assert.Equal(this.T(), "entity1", entity[0].Name())
	assert.Equal(this.T(), EntityID(2), entity[1].EntityID())
	assert.Equal(this.T(), "entity2", entity[1].Name())
}
