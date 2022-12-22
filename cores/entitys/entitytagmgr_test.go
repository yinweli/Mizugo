package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEntityTagmgr(t *testing.T) {
	suite.Run(t, new(SuiteEntityTagmgr))
}

type SuiteEntityTagmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEntityTagmgr) SetupSuite() {
	this.Change("test-entitys-entitytagmgr")
}

func (this *SuiteEntityTagmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEntityTagmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEntityTagmgr) TestNewEntityTagmgr() {
	assert.NotNil(this.T(), NewEntityTagmgr())
}

func (this *SuiteEntityTagmgr) TestEntityTagmgr() {
	target := NewEntityTagmgr()
	entity := NewEntity(EntityID(1))
	tag := []string{"tag1", "tag2", "tag3"}

	target.Add(entity, tag...)
	target.Add(entity)
	target.Add(nil)

	for _, itor := range tag {
		assert.ElementsMatch(this.T(), []*Entity{entity}, target.Get(itor))
	} // for

	assert.ElementsMatch(this.T(), tag, target.Tag(entity))
	assert.ElementsMatch(this.T(), tag, entity.Tag())

	target.Del(entity, tag...)
	target.Del(entity)
	target.Del(nil)

	for _, itor := range tag {
		assert.ElementsMatch(this.T(), []*Entity{}, target.Get(itor))
	} // for

	assert.ElementsMatch(this.T(), []string{}, target.Tag(entity))
	assert.ElementsMatch(this.T(), []string{}, entity.Tag())
}

func (this *SuiteEntityTagmgr) TestFind() {
	target := NewEntityTagmgr()
	assert.NotNil(this.T(), target.find("tag1"))
	assert.NotNil(this.T(), target.find("tag2"))
	assert.NotNil(this.T(), target.find("tag3"))
}

func (this *SuiteEntityTagmgr) TestNewEntityTag() {
	assert.NotNil(this.T(), newEntityTag())
}

func (this *SuiteEntityTagmgr) TestEntityTag() {
	target := newEntityTag()
	tag := []string{"tag1", "tag2", "tag3"}
	target.Add(tag...)
	assert.ElementsMatch(this.T(), tag, target.Tag())
	target.Del(tag...)
	assert.ElementsMatch(this.T(), []string{}, target.Tag())
}
