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
	this.Change("test-events-eventmgr")
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
	target := NewEventmgr(10)
	target.Initialize()
	target.Finalize()
}

func (this *SuiteEventmgr) TestPubOnce() {
	target := NewEventmgr(10)
	target.Initialize()

	name := "event once"
	value := "value once"
	valid := atomic.Bool{}
	target.Sub(name, func(param any) {
		valid.Store(param.(string) == value)
	})
	target.PubOnce(name, value)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubOnce(name, value)
}

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(10)
	target.Initialize()

	name := "event fixed"
	value := "value fixed"
	valid := atomic.Bool{}
	target.Sub(name, func(param any) {
		valid.Store(param.(string) == value)
	})
	target.PubFixed(name, value, testdata.Timeout)

	time.Sleep(testdata.Timeout * 5) // 多等一下讓定時事件發生
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubFixed(name, value, testdata.Timeout)
}
