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
	target := NewEntitymgr()
	assert.NotNil(this.T(), target.Add())
	target.Clear()
}

func (this *SuiteEntitymgr) TestDel() {
	target := NewEntitymgr()
	entity := target.Add()
	assert.NotNil(this.T(), entity)
	assert.NotNil(this.T(), target.Del(entity.EntityID()))
	assert.Nil(this.T(), target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Del(entity.EntityID()))
	target.Clear()
}

func (this *SuiteEntitymgr) TestClear() {
	target := NewEntitymgr()
	assert.NotNil(this.T(), target.Add())
	assert.Equal(this.T(), 1, target.Count())
	target.Clear()
	assert.Equal(this.T(), 0, target.Count())
}

func (this *SuiteEntitymgr) TestGet() {
	target := NewEntitymgr()
	entity := target.Add()
	assert.Equal(this.T(), entity, target.Get(entity.EntityID()))
	assert.Nil(this.T(), target.Get(EntityID(2)))
	target.Clear()
}

func (this *SuiteEntitymgr) TestAll() {
	target := NewEntitymgr()
	entity1 := target.Add()
	entity2 := target.Add()
	assert.ElementsMatch(this.T(), []*Entity{entity1, entity2}, target.All())
	target.Clear()
}

func (this *SuiteEntitymgr) TestCount() {
	target := NewEntitymgr()
	assert.NotNil(this.T(), target.Add())
	assert.NotNil(this.T(), target.Add())
	assert.Equal(this.T(), 2, target.Count())
	target.Clear()
}
