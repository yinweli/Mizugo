package helps

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMap(t *testing.T) {
	suite.Run(t, new(SuiteMap))
}

type SuiteMap struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMap) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-map"))
}

func (this *SuiteMap) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMap) TestMapToArray() {
	result1, result2 := MapToArray(map[int]int{
		1: 4,
		2: 5,
		3: 6,
	})
	sort.Slice(result1, func(l, r int) bool {
		return result1[l] < result1[r]
	})
	sort.Slice(result2, func(l, r int) bool {
		return result2[l] < result2[r]
	})
	assert.Equal(this.T(), []int{1, 2, 3}, result1)
	assert.Equal(this.T(), []int{4, 5, 6}, result2)
}
