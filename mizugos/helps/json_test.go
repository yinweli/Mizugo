package helps

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestJson(t *testing.T) {
	suite.Run(t, new(SuiteJson))
}

type SuiteJson struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteJson) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-json"))
}

func (this *SuiteJson) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteJson) TestJsonString() {
	this.NotEmpty(JsonString(struct{ Data int }{Data: 100}))
	this.Equal("null", JsonString(nil))
}
