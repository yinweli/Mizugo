package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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

func (this *SuiteMap) TestMapKey() {
	target := MapKey(map[int]string{
		1: "a",
		2: "b",
		3: "c",
	})
	this.ElementsMatch(target, []int{1, 2, 3})
}

func (this *SuiteMap) TestMapValue() {
	target := MapValue(map[int]string{
		1: "a",
		2: "b",
		3: "c",
	})
	this.ElementsMatch(target, []string{"a", "b", "c"})
}

func (this *SuiteMap) TestMapFlatten() {
	target := MapFlatten(map[int]string{
		1: "a",
		2: "b",
		3: "c",
	})
	this.ElementsMatch(target, []MapFlattenData[int, string]{
		{K: 1, V: "a"},
		{K: 2, V: "b"},
		{K: 3, V: "c"},
	})
}
