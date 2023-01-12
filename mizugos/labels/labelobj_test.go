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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteLabelobj) SetupSuite() {
	this.Change("test-labels-labelobj")
}

func (this *SuiteLabelobj) TearDownSuite() {
	this.Restore()
}

func (this *SuiteLabelobj) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteLabelobj) TestNewLabelobj() {
	assert.NotNil(this.T(), NewLabelobj())
}

func (this *SuiteLabelobj) TestLabelobj() {
	target := NewLabelobj()
	label := []string{"label1", "label2", "label3"}

	target.add(label...)
	assert.ElementsMatch(this.T(), label, target.label())

	target.del(label...)
	assert.ElementsMatch(this.T(), []string{}, target.label())

	target.add(label...)
	assert.ElementsMatch(this.T(), label, target.label())
	target.erase()
	assert.ElementsMatch(this.T(), []string{}, target.label())
}
