package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	// 由於需要金鑰與憑證, 因此無法測試細節
	target := NewIAPApple(&IAPAppleConfig{})
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Verify(testdata.Unknown, testdata.Unknown).Err)
	target.Finalize()
}
