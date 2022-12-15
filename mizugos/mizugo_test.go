package mizugos

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/cores/logs"
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

	assert.Nil(this.T(), Netmgr())
	assert.Nil(this.T(), Entitymgr())
	assert.Nil(this.T(), Tagmgr())
	assert.Nil(this.T(), Debug(name))
	assert.Nil(this.T(), Info(name))
	assert.Nil(this.T(), Warn(name))
	assert.Nil(this.T(), Error(name))

	go Start(name, &logs.EmptyLogger{},
		func() error {
			validInitialize.Store(true)
			return nil
		},
		func() {
			validFinalize.Store(true)
		})

	time.Sleep(time.Second)
	assert.NotNil(this.T(), Netmgr())
	assert.NotNil(this.T(), Entitymgr())
	assert.NotNil(this.T(), Tagmgr())
	assert.NotNil(this.T(), Debug(name))
	assert.NotNil(this.T(), Info(name))
	assert.NotNil(this.T(), Warn(name))
	assert.NotNil(this.T(), Error(name))

	go Close()

	time.Sleep(time.Second)
	assert.True(this.T(), validInitialize.Load())
	assert.True(this.T(), validFinalize.Load())
}

func (this *SuiteMizugo) TestTwice() {
	name := "twice"
	valid := atomic.Int64{}

	go Start(name, &logs.EmptyLogger{},
		func() error {
			valid.Add(1)
			return nil
		},
		func() {
		})
	go Start(name, &logs.EmptyLogger{},
		func() error {
			valid.Add(1)
			return nil
		},
		func() {
		})

	time.Sleep(time.Second)

	go Close()

	time.Sleep(time.Second)
	assert.Equal(this.T(), int64(1), valid.Load())
}

func (this *SuiteMizugo) TestFailed() {
	go Start("failed", &logs.EmptyLogger{},
		func() error {
			return fmt.Errorf("failed")
		},
		func() {
		})
}
