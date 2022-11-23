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

func TestEventmgr(t *testing.T) {
	suite.Run(t, new(SuiteEventmgr))
}

type SuiteEventmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteEventmgr) SetupSuite() {
	this.Change("test-eventmgr")
}

func (this *SuiteEventmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEventmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteEventmgr) TestNewEventmgr() {
	assert.NotNil(this.T(), NewEventmgr(1))
}

func (this *SuiteEventmgr) TestEventmgr() {
	target := NewEventmgr(1)
	target.Initialize()
	target.Finalize()
}

func (this *SuiteEventmgr) TestPubOnce() {
	target := NewEventmgr(10)
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

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(10)
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
