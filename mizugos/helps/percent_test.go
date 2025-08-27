package helps

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(100), target.Base())

	target = NewP1K()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(1000), target.Base())

	target = NewP10K()
	assert.NotNil(this.T(), target)
	assert.Equal(this.T(), int32(10000), target.Base())

	target = NewP100()
	assert.NotNil(this.T(), target.Set(10))
	assert.Equal(this.T(), int32(10), target.Get())
	assert.NotNil(this.T(), target.SetBase())
	assert.Equal(this.T(), int32(100), target.Get())
	assert.NotNil(this.T(), target.Add(5))
	assert.NotNil(this.T(), target.Del(5))
	assert.NotNil(this.T(), target.Mul(2))
	assert.NotNil(this.T(), target.Div(2))
	assert.NotNil(this.T(), target.Div(0))
	assert.Equal(this.T(), int32(100), target.Get())

	target = NewP100()
	assert.NotNil(this.T(), target.Set(10))
	assert.Equal(this.T(), 6, target.Calc(55, math.Round))
	assert.Equal(this.T(), 6, target.Calc(55, math.Ceil))
	assert.Equal(this.T(), 5, target.Calc(55, math.Floor))
	assert.Equal(this.T(), 8, target.Calc(75, math.Round))
	assert.Equal(this.T(), 8, target.Calc(75, math.Ceil))
	assert.Equal(this.T(), 7, target.Calc(75, math.Floor))
	assert.Equal(this.T(), 20000000, target.Calc(200000000, math.Round))
	assert.Equal(this.T(), int32(6), target.Calc32(55, math.Round))
	assert.Equal(this.T(), int32(6), target.Calc32(55, math.Ceil))
	assert.Equal(this.T(), int32(5), target.Calc32(55, math.Floor))
	assert.Equal(this.T(), int32(8), target.Calc32(75, math.Round))
	assert.Equal(this.T(), int32(8), target.Calc32(75, math.Ceil))
	assert.Equal(this.T(), int32(7), target.Calc32(75, math.Floor))
	assert.Equal(this.T(), int32(20000000), target.Calc32(200000000, math.Round))
	assert.Equal(this.T(), int64(6), target.Calc64(55, math.Round))
	assert.Equal(this.T(), int64(6), target.Calc64(55, math.Ceil))
	assert.Equal(this.T(), int64(5), target.Calc64(55, math.Floor))
	assert.Equal(this.T(), int64(8), target.Calc64(75, math.Round))
	assert.Equal(this.T(), int64(8), target.Calc64(75, math.Ceil))
	assert.Equal(this.T(), int64(7), target.Calc64(75, math.Floor))
	assert.Equal(this.T(), int64(20000000), target.Calc64(200000000, math.Round))

	target = NewPercent(0)
	assert.Zero(this.T(), target.calc(float64(1), math.Round))

	target = NewPercent(1)
	assert.Zero(this.T(), target.calc(float64(1), nil))
}
