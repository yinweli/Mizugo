package mizugos

import (
	"testing"
	"time"

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
	time.Sleep(trials.Timeout)
	assert.NotNil(this.T(), Config)
	assert.NotNil(this.T(), Metrics)
	assert.NotNil(this.T(), Logger)
	assert.NotNil(this.T(), Network)
	assert.NotNil(this.T(), Redmo)
	assert.NotNil(this.T(), Entity)
	assert.NotNil(this.T(), Label)
	assert.NotNil(this.T(), Pool)
	Stop()
	time.Sleep(trials.Timeout)
	assert.Nil(this.T(), Config)
	assert.Nil(this.T(), Metrics)
	assert.Nil(this.T(), Logger)
	assert.Nil(this.T(), Network)
	assert.Nil(this.T(), Redmo)
	assert.Nil(this.T(), Entity)
	assert.Nil(this.T(), Label)
	assert.Nil(this.T(), Pool)
}
