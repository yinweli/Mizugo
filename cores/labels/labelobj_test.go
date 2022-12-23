package labels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestLabelobj(t *testing.T) {
	suite.Run(t, new(SuiteLabelobj))
}

type SuiteLabelobj struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteLabelobj) SetupSuite() {
	this.Change("test-labels-labelobj")
}

func (this *SuiteLabelobj) TearDownSuite() {
	this.Restore()
}

func (this *SuiteLabelobj) TearDownTest() {
	goleak.VerifyNone(this.T())
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
