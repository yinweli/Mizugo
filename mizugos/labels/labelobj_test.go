package labels

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestLabelobj(t *testing.T) {
	suite.Run(t, new(SuiteLabelobj))
}

type SuiteLabelobj struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteLabelobj) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-labels-labelobj"))
}

func (this *SuiteLabelobj) TearDownSuite() {
	trials.Restore(this.Catalog)
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
