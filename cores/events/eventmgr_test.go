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
	timeout time.Duration
}

func (this *SuiteEventmgr) SetupSuite() {
	this.Change("test-events-eventmgr")
	this.timeout = time.Second
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
		valid.Store(param.(string) == "pubonce")
	})
	target.PubOnce("event", "pubonce")
	time.Sleep(this.timeout)
	assert.True(this.T(), valid.Load())
}

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(10)
	target.Initialize()
	defer target.Finalize()

	valid := atomic.Bool{}
	target.Sub("event", func(param any) {
		valid.Store(param.(string) == "pubfixed")
	})
	fixed := target.PubFixed("event", "pubfixed", time.Millisecond)
	defer fixed.Stop()
	time.Sleep(this.timeout)
	assert.True(this.T(), valid.Load())
}
