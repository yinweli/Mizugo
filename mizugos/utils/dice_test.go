package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestDice(t *testing.T) {
	suite.Run(t, new(SuiteDice))
}

type SuiteDice struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteDice) SetupSuite() {
	this.Env = testdata.EnvSetup("test-utils-dice")
}

func (this *SuiteDice) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteDice) TearDownTest() {
	testdata.Leak(this.T(), true)
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
	assert.NotNil(this.T(), target.Rand())
	assert.NotNil(this.T(), target.Randn(10))
	assert.Equal(this.T(), int64(10), target.Max())

	target.Clear()
	assert.False(this.T(), target.Valid())
	assert.Nil(this.T(), target.Rand())
	assert.Nil(this.T(), target.Randn(10))
}

func (this *SuiteDice) TestDetect() {
	target := Dice{}
	_ = target.Fill([]any{true, false}, []int64{10000, 0})

	detect := NewDiceDetect()

	for i := 0; i < testdata.TestCount; i++ {
		detect.Add(target.Rand(), 1)
	} // for

	assert.True(this.T(), detect.Check(true, testdata.TestCount, 1, 1))
	assert.True(this.T(), detect.Check(false, testdata.TestCount, 0, 0))

	target = Dice{}
	_ = target.Fill([]any{true, false}, []int64{0, 10000})

	detect = NewDiceDetect()

	for i := 0; i < testdata.TestCount; i++ {
		detect.Add(target.Rand(), 1)
	} // for

	assert.True(this.T(), detect.Check(true, testdata.TestCount, 0, 0))
	assert.True(this.T(), detect.Check(false, testdata.TestCount, 1, 1))

	target = Dice{}
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
