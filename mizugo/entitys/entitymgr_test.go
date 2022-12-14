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

func (this *SuiteEntitymgr) TestAdd() {
	entity := NewEntity(EntityID(1))
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(entity))
	assert.True(this.T(), entity.Enable())
	assert.NotNil(this.T(), target.Get(entity.EntityID()))
	assert.NotNil(this.T(), target.Add(entity))

	target.Clear()
}

func (this *SuiteEntitymgr) TestDel() {
	entity := NewEntity(EntityID(1))
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(entity))
	assert.Equal(this.T(), entity, target.Del(entity.EntityID()))
	assert.False(this.T(), entity.Enable())
	assert.Nil(this.T(), target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Del(entity.EntityID()))

	target.Clear()
}

func (this *SuiteEntitymgr) TestClear() {
	entity := NewEntity(EntityID(1))
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(entity))
	assert.Equal(this.T(), 1, target.Count())
	target.Clear()
	assert.Equal(this.T(), 0, target.Count())
}

func (this *SuiteEntitymgr) TestGet() {
	entity := NewEntity(EntityID(1))
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(entity))
	assert.Equal(this.T(), entity, target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Get(EntityID(2)))

	target.Clear()
}

func (this *SuiteEntitymgr) TestAll() {
	entity1 := NewEntity(EntityID(1))
	entity2 := NewEntity(EntityID(2))
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(entity1))
	assert.Nil(this.T(), target.Add(entity2))
	assert.ElementsMatch(this.T(), []*Entity{entity1, entity2}, target.All())

	target.Clear()
}

func (this *SuiteEntitymgr) TestCount() {
	target := NewEntitymgr()

	assert.Nil(this.T(), target.Add(NewEntity(EntityID(1))))
	assert.Nil(this.T(), target.Add(NewEntity(EntityID(2))))
	assert.Equal(this.T(), 2, target.Count())

	target.Clear()
}
