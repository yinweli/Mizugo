package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTagmgr(t *testing.T) {
	suite.Run(t, new(SuiteTagmgr))
}

type SuiteTagmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteTagmgr) SetupSuite() {
	this.Change("test-tags-tagmgr")
}

func (this *SuiteTagmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTagmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTagmgr) TestNewTagmgr() {
	assert.NotNil(this.T(), NewTagmgr())
}

func (this *SuiteTagmgr) TestTagmgr() {
	tag := []string{"tag1", "tag2", "tag3"}
	value := "value"
	target := NewTagmgr()

	target.Add(value, tag...)

	for _, itor := range tag {
		assert.ElementsMatch(this.T(), []any{value}, target.Get(itor))
	} // for

	assert.ElementsMatch(this.T(), tag, target.Tag(value))

	target.Del(value, tag...)

	for _, itor := range tag {
		assert.ElementsMatch(this.T(), []any{}, target.Get(itor))
	} // for

	assert.ElementsMatch(this.T(), []string{}, target.Tag(value))
}

func (this *SuiteTagmgr) TestFind() {
	target := NewTagmgr()
	assert.NotNil(this.T(), target.find("tag1"))
	assert.NotNil(this.T(), target.find("tag2"))
	assert.NotNil(this.T(), target.find("tag3"))
}
