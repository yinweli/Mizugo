package mizugos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMizugo(t *testing.T) {
	suite.Run(t, new(SuiteMizugo))
}

type SuiteMizugo struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteMizugo) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-mizugos-mizugo"))
}

func (this *SuiteMizugo) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteMizugo) TestMizugo() {
	Start()
	trials.WaitTimeout()
	assert.NotNil(this.T(), Config)
	assert.NotNil(this.T(), Logger)
	assert.NotNil(this.T(), Network)
	assert.NotNil(this.T(), Redmo)
	assert.NotNil(this.T(), Entity)
	assert.NotNil(this.T(), Pool)
	Stop()
	trials.WaitTimeout()
	assert.Nil(this.T(), Config)
	assert.Nil(this.T(), Logger)
	assert.Nil(this.T(), Network)
	assert.Nil(this.T(), Redmo)
	assert.Nil(this.T(), Entity)
	assert.Nil(this.T(), Pool)
}
