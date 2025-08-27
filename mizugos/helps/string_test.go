package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestString(t *testing.T) {
	suite.Run(t, new(SuiteString))
}

type SuiteString struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteString) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-string"))
}

func (this *SuiteString) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteString) TestString() {
	assert.Equal(this.T(), 30, StringDisplayLength("Hello, こんにちは, 안녕하세요!"))
	assert.NotEmpty(this.T(), StrPercentage(1, 100))
}
