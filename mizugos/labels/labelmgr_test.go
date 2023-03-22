package labels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestLabelmgr(t *testing.T) {
	suite.Run(t, new(SuiteLabelmgr))
}

type SuiteLabelmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteLabelmgr) SetupSuite() {
	this.TBegin("test-labels-labelmgr", "")
}

func (this *SuiteLabelmgr) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteLabelmgr) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteLabelmgr) TestNewLabelmgr() {
	assert.NotNil(this.T(), NewLabelmgr())
}

func (this *SuiteLabelmgr) TestLabelmgr() {
	target := NewLabelmgr()
	label := []string{"label1", "label2", "label3"}
	obj := &labelTester{
		Labelobj: NewLabelobj(),
	}

	target.Add(obj, label...)
	assert.Equal(this.T(), []any{obj}, target.Get(label[0]))
	assert.ElementsMatch(this.T(), label, target.Label(obj))

	target.Del(obj, label...)
	assert.Equal(this.T(), []any{}, target.Get(label[0]))
	assert.ElementsMatch(this.T(), []string{}, target.Label(obj))
	assert.ElementsMatch(this.T(), []string{}, target.Label(nil))

	target.Add(obj, label...)
	assert.Equal(this.T(), []any{obj}, target.Get(label[0]))
	assert.ElementsMatch(this.T(), label, target.Label(obj))
	target.Erase(obj)
	assert.Equal(this.T(), []any{}, target.Get(label[0]))
	assert.ElementsMatch(this.T(), []string{}, target.Label(obj))
}

type labelTester struct {
	*Labelobj
}
