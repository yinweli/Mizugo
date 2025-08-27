package helps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func (this *SuiteRand) TestRand() {
	RandSeed(0)
	RandSeedTime()
	assert.NotNil(this.T(), RandSource())
	fmt.Println(RandInt())
	value := RandIntn(-5, 5)
	assert.True(this.T(), value >= -5 && value <= 5)
	fmt.Println(RandInt32())
	value32 := RandInt32n(-5, 5)
	assert.True(this.T(), value32 >= -5 && value32 <= 5)
	fmt.Println(RandInt64())
	value64 := RandInt64n(-5, 5)
	assert.True(this.T(), value64 >= -5 && value64 <= 5)
	fmt.Println(RandReal64())
	value64 = RandReal64n(-5, 5)
	assert.True(this.T(), value64 >= -5 && value64 <= 5)
	values := RandString(32, StrNumberAlpha)
	assert.NotNil(this.T(), values)
	assert.Len(this.T(), values, 32)
	fmt.Println(values)
	values = RandString(64, StrNumberAlpha)
	assert.NotNil(this.T(), values)
	assert.Len(this.T(), values, 64)
	fmt.Println(values)
	values = RandStringDefault()
	assert.NotNil(this.T(), values)
	assert.Len(this.T(), values, 10)
	fmt.Println(values)
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
