package entitys

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestEntitymgr(t *testing.T) {
	suite.Run(t, new(SuiteEntitymgr))
}

type SuiteEntitymgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteEntitymgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-entitys-entitymgr"))
}

func (this *SuiteEntitymgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteEntitymgr) TestEntitymgr() {
	target := NewEntitymgr()
	this.NotNil(target)
	this.NotNil(target.Add())
	target.Clear()

	target = NewEntitymgr()
	entity := target.Add()
	this.NotNil(entity)
	this.NotNil(target.Del(entity.EntityID()))
	this.Nil(target.Get(entity.EntityID()))
	this.Nil(target.Del(entity.EntityID()))
	target.Clear()

	target = NewEntitymgr()
	this.NotNil(target.Add())
	this.Equal(1, target.Count())
	target.Clear()
	this.Equal(0, target.Count())

	target = NewEntitymgr()
	entity = target.Add()
	this.Equal(entity, target.Get(entity.EntityID()))
	this.Nil(target.Get(EntityID(2)))
	target.Clear()

	target = NewEntitymgr()
	entity1 := target.Add()
	entity2 := target.Add()
	this.ElementsMatch([]*Entity{entity1, entity2}, target.All())
	target.Clear()

	target = NewEntitymgr()
	this.NotNil(target.Add())
	this.NotNil(target.Add())
	this.Equal(2, target.Count())
	target.Clear()
}
