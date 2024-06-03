package iaps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestIAPOneStore(t *testing.T) {
	suite.Run(t, new(SuiteIAPOneStore))
}

type SuiteIAPOneStore struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteIAPOneStore) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-iapOneStore"))
}

func (this *SuiteIAPOneStore) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteIAPOneStore) TestIAPOneStore() {
	// 由於需要金鑰與憑證, 因此無法測試細節
	target := NewIAPOneStore(&IAPOneStoreConfig{})
	assert.NotNil(this.T(), target)
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Verify(testdata.Unknown, testdata.Unknown))
	target.Finalize()
}
