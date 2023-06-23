package labels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestLabelobj(t *testing.T) {
	suite.Run(t, new(SuiteLabelobj))
}

type SuiteLabelobj struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteLabelobj) SetupSuite() {
	this.Env = testdata.EnvSetup("test-labels-labelobj")
}

func (this *SuiteLabelobj) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteLabelobj) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteLabelobj) TestLabelobj() {
	label := []string{"label1", "label2", "label3"}
	target := NewLabelobj()
	assert.NotNil(this.T(), target)

	target.add(label...)
	assert.ElementsMatch(this.T(), label, target.label())

	target.del(label...)
	assert.ElementsMatch(this.T(), []string{}, target.label())

	target.add(label...)
	assert.ElementsMatch(this.T(), label, target.label())
	target.erase()
	assert.ElementsMatch(this.T(), []string{}, target.label())
}
