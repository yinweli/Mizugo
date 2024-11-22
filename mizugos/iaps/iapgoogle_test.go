package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestIAPGoogle(t *testing.T) {
	suite.Run(t, new(SuiteIAPGoogle))
}

type SuiteIAPGoogle struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPGoogle) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-iapGoogle"))
}

func (this *SuiteIAPGoogle) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPGoogle) TestIAPGoogle() {
	// 由於需要金鑰與憑證, 因此無法測試細節
	target := NewIAPGoogle(&IAPGoogleConfig{})
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Verify(testdata.Unknown, testdata.Unknown).Err)
	assert.Panics(this.T(), func() {
		target.Finalize()
	})
}
