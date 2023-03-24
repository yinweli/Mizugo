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
	assert.NotNil(this.T(), Metricsmgr())
	assert.NotNil(this.T(), Logmgr())
	assert.NotNil(this.T(), Netmgr())
	assert.NotNil(this.T(), Redmomgr())
	assert.NotNil(this.T(), Entitymgr())
	assert.NotNil(this.T(), Labelmgr())
	assert.NotNil(this.T(), Poolmgr())
	assert.NotNil(this.T(), Debug("", ""))
	assert.NotNil(this.T(), Info("", ""))
	assert.NotNil(this.T(), Warn("", ""))
	assert.NotNil(this.T(), Error("", ""))
	time.Sleep(testdata.Timeout)
	Stop()
}
