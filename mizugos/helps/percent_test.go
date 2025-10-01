package helps

import (
	"math"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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

func (this *SuitePercent) TestPercent() {
	target := NewP100()
	this.NotNil(target)
	this.Equal(PercentRatio100, target.Base())

	target = NewP1K()
	this.NotNil(target)
	this.Equal(PercentRatio1K, target.Base())

	target = NewP10K()
	this.NotNil(target)
	this.Equal(PercentRatio10K, target.Base())
}

func (this *SuitePercent) TestSet() {
	target := NewP100()
	this.NotNil(target.Set(1))
	this.Equal(int32(1), target.Get())
	this.NotNil(target.SetBase())
	this.Equal(PercentRatio100, target.Get())
}

func (this *SuitePercent) TestAdd() {
	target := NewP100()
	this.NotNil(target.Add(1))
	this.Equal(int32(1), target.Get())
}

func (this *SuitePercent) TestSub() {
	target := NewP100()
	this.NotNil(target.Set(1))
	this.NotNil(target.Sub(1))
	this.Zero(target.Get())
}

func (this *SuitePercent) TestMul() {
	target := NewP100()
	this.NotNil(target.Set(1))
	this.NotNil(target.Mul(2))
	this.Equal(int32(2), target.Get())
}

func (this *SuitePercent) TestDiv() {
	target := NewP100()
	this.NotNil(target.Set(1))
	this.NotNil(target.Div(1))
	this.NotNil(target.Div(0))
	this.Equal(int32(1), target.Get())
}

func (this *SuitePercent) TestCalc() {
	target := NewP100()
	this.NotNil(target.Set(10))
	this.Equal(6, target.Calc(55, math.Round))
	this.Equal(6, target.Calc(55, math.Ceil))
	this.Equal(5, target.Calc(55, math.Floor))
	this.Equal(8, target.Calc(75, math.Round))
	this.Equal(8, target.Calc(75, math.Ceil))
	this.Equal(7, target.Calc(75, math.Floor))
	this.Equal(20000000, target.Calc(200000000, math.Round))
}

func (this *SuitePercent) TestCalc32() {
	target := NewP100()
	this.NotNil(target.Set(10))
	this.Equal(int32(6), target.Calc32(55, math.Round))
	this.Equal(int32(6), target.Calc32(55, math.Ceil))
	this.Equal(int32(5), target.Calc32(55, math.Floor))
	this.Equal(int32(8), target.Calc32(75, math.Round))
	this.Equal(int32(8), target.Calc32(75, math.Ceil))
	this.Equal(int32(7), target.Calc32(75, math.Floor))
	this.Equal(int32(20000000), target.Calc32(200000000, math.Round))
}

func (this *SuitePercent) TestCalc64() {
	target := NewP100()
	this.NotNil(target.Set(10))
	this.Equal(int64(6), target.Calc64(55, math.Round))
	this.Equal(int64(6), target.Calc64(55, math.Ceil))
	this.Equal(int64(5), target.Calc64(55, math.Floor))
	this.Equal(int64(8), target.Calc64(75, math.Round))
	this.Equal(int64(8), target.Calc64(75, math.Ceil))
	this.Equal(int64(7), target.Calc64(75, math.Floor))
	this.Equal(int64(20000000), target.Calc64(200000000, math.Round))
}

func (this *SuitePercent) TestCalculate() {
	target := NewPercent(1)
	this.Zero(target.calculate(1, math.Round))

	target = NewPercent(0)
	this.Zero(target.calculate(1, math.Round))

	target = NewPercent(1)
	this.Zero(target.calculate(1, nil))
}
