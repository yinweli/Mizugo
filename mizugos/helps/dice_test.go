package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestDice(t *testing.T) {
	suite.Run(t, new(SuiteDice))
}

type SuiteDice struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteDice) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-dice"))
}

func (this *SuiteDice) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteDice) TestDice() {
	target := NewDice()
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.One("1", 1))
	assert.Nil(this.T(), target.One("2", 0))
	assert.NotNil(this.T(), target.One("3", -1))
	assert.Nil(this.T(), target.Fill([]any{"4", "5", "6"}, []int64{1, 1, 1}))
	assert.Nil(this.T(), target.Fill([]any{"7", "8", "9"}, []int64{0, 0, 0}))
	assert.NotNil(this.T(), target.Fill([]any{"10", "11", "12"}, []int64{-1, -1, -1}))
	assert.NotNil(this.T(), target.Fill([]any{"13", "14", "15"}, []int64{1}))
	assert.Nil(this.T(), target.Complete("16", 10))
	assert.Nil(this.T(), target.Complete("17", 1))
	assert.True(this.T(), target.Valid())
	assert.Equal(this.T(), int64(10), target.Max())
	assert.NotNil(this.T(), target.Rand())
	assert.NotNil(this.T(), target.Randn(10))
	assert.NotNil(this.T(), target.RandOnce())

	target.Clear()
	assert.False(this.T(), target.Valid())
	assert.Nil(this.T(), target.Rand())
	assert.Nil(this.T(), target.Randn(10))
	assert.Nil(this.T(), target.RandOnce())
}

func (this *SuiteDice) TestDiceRand() {
	target := NewDice()
	_ = target.Fill([]any{true, false}, []int64{10000, 0})

	detect := NewDiceDetect()

	for i := 0; i < testdata.TestCount; i++ {
		detect.Add(target.Rand(), 1)
	} // for

	assert.True(this.T(), detect.Check(true, testdata.TestCount, 1, 1))
	assert.True(this.T(), detect.Check(false, testdata.TestCount, 0, 0))

	target = NewDice()
	_ = target.Fill([]any{true, false}, []int64{0, 10000})

	detect = NewDiceDetect()

	for i := 0; i < testdata.TestCount; i++ {
		detect.Add(target.Rand(), 1)
	} // for

	assert.True(this.T(), detect.Check(true, testdata.TestCount, 0, 0))
	assert.True(this.T(), detect.Check(false, testdata.TestCount, 1, 1))

	target = NewDice()
	_ = target.Fill([]any{1, 2, 3, 4, 5}, []int64{500, 1000, 1500, 2000, 2500})
	_ = target.Complete(0, 10000)
	detect = NewDiceDetect()

	for i := 0; i < testdata.TestCount; i++ {
		detect.Add(target.Rand(), 1)
	} // for

	assert.True(this.T(), detect.Check(1, testdata.TestCount, 0.04, 0.06))
	assert.True(this.T(), detect.Check(2, testdata.TestCount, 0.09, 0.11))
	assert.True(this.T(), detect.Check(3, testdata.TestCount, 0.14, 0.16))
	assert.True(this.T(), detect.Check(4, testdata.TestCount, 0.19, 0.21))
	assert.True(this.T(), detect.Check(5, testdata.TestCount, 0.24, 0.26))
	assert.True(this.T(), detect.Check(0, testdata.TestCount, 0.24, 0.26))
}

func (this *SuiteDice) TestDiceRandOnce() {
	target := NewDice()
	payload := []any{1, 2, 3, 4}
	weight := []int64{10, 10, 10, 10}
	_ = target.Fill(payload, weight)
	assert.Contains(this.T(), payload, target.RandOnce())
	assert.Contains(this.T(), payload, target.RandOnce())
	assert.Contains(this.T(), payload, target.RandOnce())
	assert.Contains(this.T(), payload, target.RandOnce())
	assert.False(this.T(), target.Valid())
}

func (this *SuiteDice) TestDiceDetect() {
	target := NewDiceDetect()
	assert.NotNil(this.T(), target)
	target.Add(1, 1)
	target.Add(2, 2)
	target.Add(3, 3)
	target.Add(4, 4)
	assert.Equal(this.T(), 0.1, target.Ratio(1, 10))
	assert.Equal(this.T(), 0.2, target.Ratio(2, 10))
	assert.Equal(this.T(), 0.3, target.Ratio(3, 10))
	assert.Equal(this.T(), 0.4, target.Ratio(4, 10))
	assert.True(this.T(), target.Check(4, 10, 0.1, 0.4))
	assert.False(this.T(), target.Check(4, 10, 0.1, 0.1))
}
