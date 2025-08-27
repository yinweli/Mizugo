package cryptos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestDefine(t *testing.T) {
	suite.Run(t, new(SuiteDefine))
}

type SuiteDefine struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteDefine) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-cryptos-define"))
}

func (this *SuiteDefine) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteDefine) TestRandDesKey() {
	key1 := RandDesKey()
	assert.NotNil(this.T(), key1)
	assert.Len(this.T(), key1, DesKeySize)

	key2 := RandDesKeyString()
	assert.NotNil(this.T(), key2)
	assert.Len(this.T(), key2, DesKeySize)
}
