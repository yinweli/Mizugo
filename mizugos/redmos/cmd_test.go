package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMX(t *testing.T) {
	suite.Run(t, new(SuiteMX))
}

type SuiteMX struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMX) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-mx"))
}

func (this *SuiteMX) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMX) TestSave() {
	target := NewSave()
	assert.NotNil(this.T(), target)
	target.SetSave()
	assert.True(this.T(), target.GetSave())
}
