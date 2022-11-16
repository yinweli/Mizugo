package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

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

func (this *SuiteEntityan) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntityan) TestNewEntityan() {
	assert.NotNil(this.T(), NewEntityan())
}

func (this *SuiteEntityan) TestClear() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.Equal(this.T(), 1, target.Count())
	target.Clear()
	assert.Equal(this.T(), 0, target.Count())
}

func (this *SuiteEntityan) TestAdd() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.NotNil(this.T(), target.Get(EntityID(1)))

	assert.NotNil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
}

func (this *SuiteEntityan) TestDel() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	entity := target.Del(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())
	assert.Nil(this.T(), target.Get(EntityID(1)))

	assert.Nil(this.T(), target.Del(EntityID(2)))
}

func (this *SuiteEntityan) TestGet() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	entity := target.Get(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	assert.Nil(this.T(), target.Get(EntityID(2)))
}

func (this *SuiteEntityan) TestAll() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity1")))
	assert.Nil(this.T(), target.Add(NewEntity(EntityID(2), "entity2")))
	entity := target.All()
	assert.Len(this.T(), entity, 2)
	assert.Equal(this.T(), EntityID(1), entity[0].EntityID())
	assert.Equal(this.T(), "entity1", entity[0].Name())
	assert.Equal(this.T(), EntityID(2), entity[1].EntityID())
	assert.Equal(this.T(), "entity2", entity[1].Name())
}

func (this *SuiteEntityan) TestCount() {
	target := NewEntityan()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity1")))
	assert.Nil(this.T(), target.Add(NewEntity(EntityID(2), "entity2")))
	assert.Equal(this.T(), 2, target.Count())
}
