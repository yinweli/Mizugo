package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestMath(t *testing.T) {
	suite.Run(t, new(SuiteMath))
}

type SuiteMath struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMath) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-math"))
}

func (this *SuiteMath) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMath) TestMin() {
	this.Equal(1, Min([]int{5, 4, 3, 2, 1}, 0))
	this.Equal(0, Min([]int{}, 0))
}

func (this *SuiteMath) TestMax() {
	this.Equal(5, Max([]int{5, 4, 3, 2, 1}, 0))
	this.Equal(0, Max([]int{}, 0))
}
