package events

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEventan(t *testing.T) {
	suite.Run(t, new(SuiteEventan))
}

type SuiteEventan struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEventan) SetupSuite() {
	this.Change("test-eventan")
}

func (this *SuiteEventan) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEventan) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEventan) TestNewEventan() {
	assert.NotNil(this.T(), NewEventan(1))
}

func (this *SuiteEventan) TestEventan() {
	target := NewEventan(1)
	target.Initialize()
	target.Finalize()
}

func (this *SuiteEventan) TestPubOnce() {
	target := NewEventan(10)
	target.Initialize()
	defer target.Finalize()

	valid := atomic.Bool{}
	target.Sub("event", func(param any) {
		if param.(string) == "pubonce" {
			valid.Store(true)
		} // if
	})
	target.PubOnce("event", "pubonce")
	time.Sleep(time.Millisecond * 10)
	assert.True(this.T(), valid.Load())
}

func (this *SuiteEventan) TestPubFixed() {
	target := NewEventan(10)
	target.Initialize()
	defer target.Finalize()

	valid := atomic.Int64{}
	target.Sub("event", func(param any) {
		if param.(string) == "pubfixed" {
			valid.Add(1)
		} // if
	})
	fixed := target.PubFixed("event", "pubfixed", time.Millisecond)
	defer fixed.Stop()
	time.Sleep(time.Millisecond * 100)
	assert.Greater(this.T(), valid.Load(), int64(0))
}
