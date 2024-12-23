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

func (this *SuiteMap) TestMapFlatten() {
	target := MapFlatten(map[int]string{
		1: "a",
		2: "b",
		3: "c",
	})
	sort.Slice(target, func(l, r int) bool {
		return target[l].K < target[r].K
	})
	assert.Equal(this.T(), []MapFlattenData[int, string]{
		{K: 1, V: "a"},
		{K: 2, V: "b"},
		{K: 3, V: "c"},
	}, target)
}
