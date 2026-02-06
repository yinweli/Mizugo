package redmos

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestSave(t *testing.T) {
	suite.Run(t, new(SuiteSave))
}

type SuiteSave struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteSave) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-redmos-save"))
}

func (this *SuiteSave) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteSave) TestSave() {
	target := NewSave()
	this.NotNil(target)
	target.SetSave()
	this.True(target.GetSave())
	target.ClrSave()
	this.False(target.GetSave())
}
