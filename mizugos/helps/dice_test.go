package helps

import (
	"testing"

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
	this.NotNil(target)
	target.Clear()
	this.False(target.Valid())
}

func (this *SuiteDice) TestOne() {
	target := NewDice()
	this.Nil(target.One("a", 1))
	this.Nil(target.One("b", 0))
	this.NotNil(target.One("c", -1))
	this.True(target.Valid())
	this.Equal(int64(1), target.Max())
}

func (this *SuiteDice) TestFill() {
	target := NewDice()
	this.Nil(target.Fill([]any{"a", "b", "c"}, []int64{1, 1, 1}))
	this.Nil(target.Fill([]any{"d", "e", "f"}, []int64{0, 0, 0}))
	this.NotNil(target.Fill([]any{"g", "h", "i"}, []int64{-1, -1, -1}))
	this.NotNil(target.Fill([]any{"j", "k", "l"}, []int64{1}))
	this.True(target.Valid())
	this.Equal(int64(3), target.Max())
}

func (this *SuiteDice) TestComplete() {
	target := NewDice()
	this.Nil(target.One("a", 1))
	this.Nil(target.Complete("b", 3))
	this.Nil(target.Complete("c", 0))
	this.True(target.Valid())
	this.Equal(int64(3), target.Max())
}

func (this *SuiteDice) TestRand() {
	target := NewDice()
	_ = target.Fill([]any{"a", "b"}, []int64{10000, 0})
	tester := newDiceTester()

	for i := 0; i < testdata.TestCount; i++ {
		tester.Add(target.Rand(), 1)
	} // for

	this.True(tester.Check("a", testdata.TestCount, 1, 1))
	this.True(tester.Check("b", testdata.TestCount, 0, 0))

	target = NewDice()
	_ = target.Fill([]any{"a", "b"}, []int64{0, 10000})
	tester = newDiceTester()

	for i := 0; i < testdata.TestCount; i++ {
		tester.Add(target.Rand(), 1)
	} // for

	this.True(tester.Check("a", testdata.TestCount, 0, 0))
	this.True(tester.Check("b", testdata.TestCount, 1, 1))

	target = NewDice()
	_ = target.Fill([]any{"a", "b", "c", "d", "e"}, []int64{500, 1000, 1500, 2000, 2500})
	_ = target.Complete("f", 10000)
	tester = newDiceTester()

	for i := 0; i < testdata.TestCount; i++ {
		tester.Add(target.Rand(), 1)
	} // for

	this.True(tester.Check("a", testdata.TestCount, 0.04, 0.06))
	this.True(tester.Check("b", testdata.TestCount, 0.09, 0.11))
	this.True(tester.Check("c", testdata.TestCount, 0.14, 0.16))
	this.True(tester.Check("d", testdata.TestCount, 0.19, 0.21))
	this.True(tester.Check("e", testdata.TestCount, 0.24, 0.26))
	this.True(tester.Check("f", testdata.TestCount, 0.24, 0.26))

	target = NewDice()
	this.Nil(target.Rand())
}

func (this *SuiteDice) TestRandOnce() {
	target := NewDice()
	payload := []any{"a", "b", "c", "d"}
	_ = target.Fill(payload, []int64{10, 10, 10, 10})
	this.Contains(payload, target.RandOnce())
	this.Contains(payload, target.RandOnce())
	this.Contains(payload, target.RandOnce())
	this.Contains(payload, target.RandOnce())
	this.Nil(target.RandOnce())
	this.False(target.Valid())
}

func newDiceTester() *testDice {
	return &testDice{
		data: map[any]int{},
	}
}

type testDice struct {
	data map[any]int
}

func (this *testDice) Add(key any, count int) {
	this.data[key] += count
}

func (this *testDice) Check(key any, total int, minimum, maximum float64) bool {
	ratio := float64(this.data[key]) / float64(total)
	return ratio >= minimum && ratio <= maximum
}
