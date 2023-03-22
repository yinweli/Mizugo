package mizugos

import (
	"fmt"
	"sync/atomic"
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
	name := "mizugo"
	tester := &mizugoTester{}

	go Start(name, tester.initialize, tester.finalize, tester.crashlize)
	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), name, Name())
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
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), tester.validInit())
	assert.True(this.T(), tester.validFinal())

	go Start(name, func() error {
		return fmt.Errorf("failed")
	}, nil, nil)
	time.Sleep(testdata.Timeout)
	Stop()

	go Start(name, nil, nil, nil)
	time.Sleep(testdata.Timeout)
	Stop()
}

type mizugoTester struct {
	countInit atomic.Int64
	countFin  atomic.Int64
}

func (this *mizugoTester) initialize() error {
	this.countInit.Add(1)
	return nil
}

func (this *mizugoTester) finalize() {
	this.countFin.Add(1)
}

func (this *mizugoTester) crashlize(_ any) {
}

func (this *mizugoTester) validInit() bool {
	return this.countInit.Load() == 1
}

func (this *mizugoTester) validFinal() bool {
	return this.countFin.Load() == 1
}
