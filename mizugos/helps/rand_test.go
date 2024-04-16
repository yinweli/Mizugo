package helps

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
