package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestPprof(t *testing.T) {
	suite.Run(t, new(SuitePprof))
}

type SuitePprof struct {
	suite.Suite
	trials.Catalog
}

func (this *SuitePprof) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-pprof"))
}

func (this *SuitePprof) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuitePprof) TestPprof() {
	target := &Pprof{}
	assert.Nil(this.T(), target.Start("test.pprof"))
	target.Stop()
}
