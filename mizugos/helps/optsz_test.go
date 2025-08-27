package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestOptsz(t *testing.T) {
	suite.Run(t, new(SuiteOptsz))
}

type SuiteOptsz struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteOptsz) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-optsz"))
}

func (this *SuiteOptsz) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteOptsz) TestOptsz() {
	target := Optsz(1)
	option := target.On("")
	assert.True(this.T(), target.Get(option))
	option = target.Off(option)
	assert.False(this.T(), target.Get(option))
}
