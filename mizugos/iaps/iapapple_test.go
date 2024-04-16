package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestIAPApple(t *testing.T) {
	suite.Run(t, new(SuiteIAPApple))
}

type SuiteIAPApple struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPApple) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-iapApple"))
}

func (this *SuiteIAPApple) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPApple) TestIAPApple() {
	target := NewIAPApple(&IAPAppleConfig{})
	assert.NotNil(this.T(), target) // 由於iap需要金鑰與憑證, 因此無法測試細節
}
