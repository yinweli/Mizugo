package mizugos

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

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
	this.Change("test-mizugos-mizugo")
}

func (this *SuiteMizugo) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMizugo) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteMizugo) TestMizugo() {
	name := "mizugo"

	tester := &mizugoTester{}
	time.Sleep(testdata.Timeout)
	go Start(name, tester.initialize, tester.finalize)
	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), name, Name())
	assert.NotNil(this.T(), Configmgr())
	assert.NotNil(this.T(), Metricsmgr())
	assert.NotNil(this.T(), Logmgr())
	assert.NotNil(this.T(), Netmgr())
	assert.NotNil(this.T(), Entitymgr())
	assert.NotNil(this.T(), Labelmgr())
	assert.NotNil(this.T(), Debug(""))
	assert.NotNil(this.T(), Info(""))
	assert.NotNil(this.T(), Warn(""))
	assert.NotNil(this.T(), Error(""))
	go Close()
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), tester.validInit())
	assert.True(this.T(), tester.validFinal())

	tester = &mizugoTester{}
	time.Sleep(testdata.Timeout)
	go Start(name, tester.initialize, tester.finalize)
	go Start(name, tester.initialize, tester.finalize)
	time.Sleep(testdata.Timeout)
	go Close()
	assert.True(this.T(), tester.validInit())
	assert.True(this.T(), tester.validFinal())

	time.Sleep(testdata.Timeout)
	go Start("nil", nil, nil)
	time.Sleep(testdata.Timeout)
	go Close()

	time.Sleep(testdata.Timeout)
	go Start("failed", func() error { return fmt.Errorf("failed") }, nil)
	time.Sleep(testdata.Timeout)
	go Close()
}

type mizugoTester struct {
	initCount  atomic.Int64
	finalCount atomic.Int64
}

func (this *mizugoTester) initialize() error {
	this.initCount.Add(1)
	return nil
}

func (this *mizugoTester) finalize() {
	this.finalCount.Add(1)
}

func (this *mizugoTester) validInit() bool {
	return this.initCount.Load() == 1
}

func (this *mizugoTester) validFinal() bool {
	return this.finalCount.Load() == 1
}
