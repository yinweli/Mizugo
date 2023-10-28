package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestIAPGoogle(t *testing.T) {
	suite.Run(t, new(SuiteIAPGoogle))
}

type SuiteIAPGoogle struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteIAPGoogle) SetupSuite() {
	this.Env = testdata.EnvSetup("test-utils-iapGoogle")
}

func (this *SuiteIAPGoogle) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteIAPGoogle) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteIAPGoogle) TestIAPGoogle() {
	target := NewIAPGoogle(&IAPGoogleConfig{})
	assert.NotNil(this.T(), target) // 由於iap需要金鑰與憑證, 因此無法測試細節
}
