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
	validInitialize := atomic.Bool{}
	validFinalize := atomic.Bool{}

	go Start(name,
		func() error {
			validInitialize.Store(true)
			return nil
		},
		func() {
			validFinalize.Store(true)
		})

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validInitialize.Load())
	assert.Equal(this.T(), name, Name())
	assert.NotNil(this.T(), Configmgr())
	assert.NotNil(this.T(), Netmgr())
	assert.NotNil(this.T(), Entitymgr())
	assert.NotNil(this.T(), Tagmgr())
	assert.NotNil(this.T(), Logmgr())
	assert.NotNil(this.T(), Debug(""))
	assert.NotNil(this.T(), Info(""))
	assert.NotNil(this.T(), Warn(""))
	assert.NotNil(this.T(), Error(""))

	go Close()

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validFinalize.Load())
}

func (this *SuiteMizugo) TestTwice() {
	valid := atomic.Int64{}

	go Start("twice",
		func() error {
			valid.Add(1)
			return nil
		},
		func() {
		})

	time.Sleep(testdata.Timeout)

	go Start("!?",
		func() error {
			valid.Add(1)
			return nil
		},
		func() {
		})

	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), int64(1), valid.Load())

	go Close()
}

func (this *SuiteMizugo) TestFailed() {
	go Start("failed",
		func() error {
			return fmt.Errorf("failed")
		},
		func() {
		})
}
