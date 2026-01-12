package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestRand(t *testing.T) {
	suite.Run(t, new(SuiteRand))
}

type SuiteRand struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteRand) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-rand"))
}

func (this *SuiteRand) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteRand) TestRandInt() {
	target := RandInt()
	this.Positive(target)
	target = RandIntn(-5, 5)
	this.True(target >= -5 && target <= 5)
}

func (this *SuiteRand) TestRandInt32() {
	target := RandInt32()
	this.Positive(target)
	target = RandInt32n(-5, 5)
	this.True(target >= -5 && target <= 5)
}

func (this *SuiteRand) TestRandInt64() {
	target := RandInt64()
	this.Positive(target)
	target = RandInt64n(-5, 5)
	this.True(target >= -5 && target <= 5)
}

func (this *SuiteRand) TestRandReal64() {
	target := RandReal64()
	this.Positive(target)
	target = RandReal64n(-5, 5)
	this.True(target >= -5 && target <= 5)
}

func (this *SuiteRand) TestRandString() {
	target := RandString(32, StrNumberAlpha)
	this.Len(target, 32)

	for _, itor := range target {
		this.Contains(StrNumberAlpha, string(itor))
	} // for

	target = RandString(64, StrNumberAlpha)
	this.Len(target, 64)

	for _, itor := range target {
		this.Contains(StrNumberAlpha, string(itor))
	} // for
}

func (this *SuiteRand) TestRandStringDefault() {
	target := RandStringDefault()
	this.Len(target, 10)

	for _, itor := range target {
		this.Contains(StrNumberAlpha, string(itor))
	} // for
}

func BenchmarkRandInt32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt32()
	} // for
}

func BenchmarkRandInt32n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt32n(0, 10000)
	} // for
}

func BenchmarkRandInt64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt64()
	} // for
}

func BenchmarkRandInt64n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandInt64n(0, 10000)
	} // for
}

func BenchmarkRandReal64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandReal64()
	} // for
}

func BenchmarkRandReal64n(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandReal64n(0, 10000)
	} // for
}

func BenchmarkRandStringDefault(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RandStringDefault()
	} // for
}
