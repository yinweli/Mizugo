package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntityan(t *testing.T) {
	suite.Run(t, new(SuiteEntityan))
}

type SuiteEntityan struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntityan) SetupSuite() {
	this.Change("test-entityan")
}

func (this *SuiteEntityan) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntityan) TestNewEntityan() {
	assert.NotNil(this.T(), NewEntityan())
}

func (this *SuiteEntityan) TestAdd() {
	target := NewEntityan()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.NotNil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
}

func (this *SuiteEntityan) TestDel() {
	target := NewEntityan()
	_ = target.Add(NewEntity(EntityID(1), "entity"))

	entity := target.Del(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	entity = target.Del(EntityID(2))
	assert.Nil(this.T(), entity)
}

func (this *SuiteEntityan) TestGet() {
	target := NewEntityan()
	_ = target.Add(NewEntity(EntityID(1), "entity"))

	entity := target.Get(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	entity = target.Get(EntityID(2))
	assert.Nil(this.T(), entity)
}

func (this *SuiteEntityan) TestAll() {
	target := NewEntityan()
	_ = target.Add(NewEntity(EntityID(1), "entity1"))
	_ = target.Add(NewEntity(EntityID(2), "entity2"))

	entity := target.All()
	assert.Len(this.T(), entity, 2)
	assert.Equal(this.T(), EntityID(1), entity[0].EntityID())
	assert.Equal(this.T(), "entity1", entity[0].Name())
	assert.Equal(this.T(), EntityID(2), entity[1].EntityID())
	assert.Equal(this.T(), "entity2", entity[1].Name())
}
