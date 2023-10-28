package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestIAPApple(t *testing.T) {
	suite.Run(t, new(SuiteIAPApple))
}

type SuiteIAPApple struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteIAPApple) SetupSuite() {
	this.Env = testdata.EnvSetup("test-utils-iapApple")
}

func (this *SuiteIAPApple) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteIAPApple) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteIAPApple) TestIAPApple() {
	target := NewIAPApple(&IAPAppleConfig{})
	assert.NotNil(this.T(), target) // 由於iap需要金鑰與憑證, 因此無法測試細節
}
