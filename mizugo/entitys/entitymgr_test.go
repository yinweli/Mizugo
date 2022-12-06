package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntitymgr(t *testing.T) {
	suite.Run(t, new(SuiteEntitymgr))
}

type SuiteEntitymgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntitymgr) SetupSuite() {
	this.Change("test-entitys-entitymgr")
}

func (this *SuiteEntitymgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntitymgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntitymgr) TestNewEntitymgr() {
	assert.NotNil(this.T(), NewEntitymgr())
}

func (this *SuiteEntitymgr) TestClear() {
	target := NewEntitymgr()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.Equal(this.T(), 1, target.Count())
	target.Clear()
	assert.Equal(this.T(), 0, target.Count())
}

func (this *SuiteEntitymgr) TestAdd() {
	target := NewEntitymgr()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	assert.NotNil(this.T(), target.Get(EntityID(1)))

	assert.NotNil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
}

func (this *SuiteEntitymgr) TestDel() {
	target := NewEntitymgr()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	entity := target.Del(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())
	assert.Nil(this.T(), target.Get(EntityID(1)))

	assert.Nil(this.T(), target.Del(EntityID(2)))
}

func (this *SuiteEntitymgr) TestGet() {
	target := NewEntitymgr()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity")))
	entity := target.Get(EntityID(1))
	assert.NotNil(this.T(), entity)
	assert.Equal(this.T(), EntityID(1), entity.EntityID())
	assert.Equal(this.T(), "entity", entity.Name())

	assert.Nil(this.T(), target.Get(EntityID(2)))
}

func (this *SuiteEntitymgr) TestAll() {
	target := NewEntitymgr()
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

func (this *SuiteEntitymgr) TestCount() {
	target := NewEntitymgr()
	defer target.Clear()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1), "entity1")))
	assert.Nil(this.T(), target.Add(NewEntity(EntityID(2), "entity2")))
	assert.Equal(this.T(), 2, target.Count())
}
