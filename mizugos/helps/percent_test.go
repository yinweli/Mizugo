package helps

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestPercent(t *testing.T) {
	suite.Run(t, new(SuitePercent))
}

type SuitePercent struct {
	suite.Suite
	trials.Catalog
}

func (this *SuitePercent) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-percent"))
}

func (this *SuitePercent) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuitePercent) TestPercent100() {
	target := NewP100()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(100), target.Base())
	assert.NotNil(this.T(), target.Set(10))
	assert.Equal(this.T(), int32(10), target.Get())
	assert.NotNil(this.T(), target.SetBase())
	assert.Equal(this.T(), int32(100), target.Get())
	assert.NotNil(this.T(), target.Add(7))
	assert.NotNil(this.T(), target.Del(5))
	assert.NotNil(this.T(), target.Mul(2))
	assert.NotNil(this.T(), target.Div(2))
	assert.NotNil(this.T(), target.Div(0))
	assert.Equal(this.T(), int32(102), target.Get())
	assert.Equal(this.T(), 56, target.Calc(55, math.Round))
	assert.Equal(this.T(), 57, target.Calc(55, math.Ceil))
	assert.Equal(this.T(), 56, target.Calc(55, math.Floor))
	assert.Equal(this.T(), 77, target.Calc(75, math.Round))
	assert.Equal(this.T(), 77, target.Calc(75, math.Ceil))
	assert.Equal(this.T(), 76, target.Calc(75, math.Floor))
	assert.Equal(this.T(), int32(56), target.Calc32(55, math.Round))
	assert.Equal(this.T(), int32(57), target.Calc32(55, math.Ceil))
	assert.Equal(this.T(), int32(56), target.Calc32(55, math.Floor))
	assert.Equal(this.T(), int32(77), target.Calc32(75, math.Round))
	assert.Equal(this.T(), int32(77), target.Calc32(75, math.Ceil))
	assert.Equal(this.T(), int32(76), target.Calc32(75, math.Floor))
	assert.Equal(this.T(), int64(56), target.Calc64(55, math.Round))
	assert.Equal(this.T(), int64(57), target.Calc64(55, math.Ceil))
	assert.Equal(this.T(), int64(56), target.Calc64(55, math.Floor))
	assert.Equal(this.T(), int64(77), target.Calc64(75, math.Round))
	assert.Equal(this.T(), int64(77), target.Calc64(75, math.Ceil))
	assert.Equal(this.T(), int64(76), target.Calc64(75, math.Floor))
}

func (this *SuitePercent) TestPercent1K() {
	target := NewP1K()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(1000), target.Base())
	assert.NotNil(this.T(), target.Set(10))
	assert.Equal(this.T(), int32(10), target.Get())
	assert.NotNil(this.T(), target.SetBase())
	assert.Equal(this.T(), int32(1000), target.Get())
	assert.NotNil(this.T(), target.Add(7))
	assert.NotNil(this.T(), target.Del(5))
	assert.NotNil(this.T(), target.Mul(2))
	assert.NotNil(this.T(), target.Div(2))
	assert.NotNil(this.T(), target.Div(0))
	assert.Equal(this.T(), int32(1002), target.Get())
	assert.Equal(this.T(), 55, target.Calc(55, math.Round))
	assert.Equal(this.T(), 56, target.Calc(55, math.Ceil))
	assert.Equal(this.T(), 55, target.Calc(55, math.Floor))
	assert.Equal(this.T(), 75, target.Calc(75, math.Round))
	assert.Equal(this.T(), 76, target.Calc(75, math.Ceil))
	assert.Equal(this.T(), 75, target.Calc(75, math.Floor))
	assert.Equal(this.T(), int32(55), target.Calc32(55, math.Round))
	assert.Equal(this.T(), int32(56), target.Calc32(55, math.Ceil))
	assert.Equal(this.T(), int32(55), target.Calc32(55, math.Floor))
	assert.Equal(this.T(), int32(75), target.Calc32(75, math.Round))
	assert.Equal(this.T(), int32(76), target.Calc32(75, math.Ceil))
	assert.Equal(this.T(), int32(75), target.Calc32(75, math.Floor))
	assert.Equal(this.T(), int64(55), target.Calc64(55, math.Round))
	assert.Equal(this.T(), int64(56), target.Calc64(55, math.Ceil))
	assert.Equal(this.T(), int64(55), target.Calc64(55, math.Floor))
	assert.Equal(this.T(), int64(75), target.Calc64(75, math.Round))
	assert.Equal(this.T(), int64(76), target.Calc64(75, math.Ceil))
	assert.Equal(this.T(), int64(75), target.Calc64(75, math.Floor))
}

func (this *SuitePercent) TestPercent10K() {
	target := NewP10K()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(10000), target.Base())
	assert.NotNil(this.T(), target.Set(10))
	assert.Equal(this.T(), int32(10), target.Get())
	assert.NotNil(this.T(), target.SetBase())
	assert.Equal(this.T(), int32(10000), target.Get())
	assert.NotNil(this.T(), target.Add(7))
	assert.NotNil(this.T(), target.Del(5))
	assert.NotNil(this.T(), target.Mul(2))
	assert.NotNil(this.T(), target.Div(2))
	assert.NotNil(this.T(), target.Div(0))
	assert.Equal(this.T(), int32(10002), target.Get())
	assert.Equal(this.T(), 55, target.Calc(55, math.Round))
	assert.Equal(this.T(), 56, target.Calc(55, math.Ceil))
	assert.Equal(this.T(), 55, target.Calc(55, math.Floor))
	assert.Equal(this.T(), 75, target.Calc(75, math.Round))
	assert.Equal(this.T(), 76, target.Calc(75, math.Ceil))
	assert.Equal(this.T(), 75, target.Calc(75, math.Floor))
	assert.Equal(this.T(), int32(55), target.Calc32(55, math.Round))
	assert.Equal(this.T(), int32(56), target.Calc32(55, math.Ceil))
	assert.Equal(this.T(), int32(55), target.Calc32(55, math.Floor))
	assert.Equal(this.T(), int32(75), target.Calc32(75, math.Round))
	assert.Equal(this.T(), int32(76), target.Calc32(75, math.Ceil))
	assert.Equal(this.T(), int32(75), target.Calc32(75, math.Floor))
	assert.Equal(this.T(), int64(55), target.Calc64(55, math.Round))
	assert.Equal(this.T(), int64(56), target.Calc64(55, math.Ceil))
	assert.Equal(this.T(), int64(55), target.Calc64(55, math.Floor))
	assert.Equal(this.T(), int64(75), target.Calc64(75, math.Round))
	assert.Equal(this.T(), int64(76), target.Calc64(75, math.Ceil))
	assert.Equal(this.T(), int64(75), target.Calc64(75, math.Floor))
}
