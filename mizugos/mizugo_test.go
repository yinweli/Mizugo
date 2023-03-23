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
	testdata.TestEnv
}

func (this *SuiteMizugo) SetupSuite() {
	this.TBegin("test-mizugos-mizugo", "")
}

func (this *SuiteMizugo) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteMizugo) TearDownTest() {
	this.TLeak(this.T(), true)
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
