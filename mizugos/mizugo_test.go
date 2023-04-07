package mizugos

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMizugo(t *testing.T) {
	suite.Run(t, new(SuiteMizugo))
}

type SuiteMizugo struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteMizugo) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-mizugos-mizugo")
}

func (this *SuiteMizugo) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteMizugo) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMizugo) TestMizugo() {
	Start()
	time.Sleep(testdata.Timeout)
	assert.NotNil(this.T(), Configmgr())
	assert.NotNil(this.T(), Varmgr())
	assert.NotNil(this.T(), Metricsmgr())
	assert.NotNil(this.T(), Logmgr())
	assert.NotNil(this.T(), Netmgr())
	assert.NotNil(this.T(), Redmomgr())
	assert.NotNil(this.T(), Entitymgr())
	assert.NotNil(this.T(), Labelmgr())
	assert.NotNil(this.T(), Poolmgr())
	Stop()
	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), Configmgr())
	assert.Nil(this.T(), Varmgr())
	assert.Nil(this.T(), Metricsmgr())
	assert.Nil(this.T(), Logmgr())
	assert.Nil(this.T(), Netmgr())
	assert.Nil(this.T(), Redmomgr())
	assert.Nil(this.T(), Entitymgr())
	assert.Nil(this.T(), Labelmgr())
	assert.Nil(this.T(), Poolmgr())
}
