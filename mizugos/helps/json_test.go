package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
	assert.NotNil(this.T(), JsonString(struct{ Data int }{Data: 100}))
	assert.NotNil(this.T(), JsonString(nil))
}
