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
	subID := ""
	subID = target.Sub(name, func(param any) {
		valid.Store(param.(string) == value)
		target.Unsub(subID)
	})

	target.PubOnce(name, value)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubOnce(name, value) // 測試在結束之後發布事件
}

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(10)
	target.Initialize()
	name := "event fixed"
	value := "value fixed"
	valid := atomic.Bool{}
	subID := ""
	subID = target.Sub(name, func(param any) {
		valid.Store(param.(string) == value)
		target.Unsub(subID)
	})

	target.PubFixed(name, value, testdata.Timeout)
	time.Sleep(testdata.Timeout * 2) // 多等一下讓定時事件發生
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubFixed(name, value, testdata.Timeout) // 測試在結束之後發布事件
}

func (this *SuiteEventmgr) TestPubsub() {
	target := newPubsub()
	name := "pubsub"
	value := "value"
	count := 0
	subID := target.sub(name, func(param any) {
		if param.(string) == value {
			count++
		} // if
	})

	target.pub(name, value)
	assert.Equal(this.T(), 1, count)

	target.unsub(subID)
	target.pub(name, value)
	assert.Equal(this.T(), 1, count)
}

func (this *SuiteEventmgr) TestSubID() {
	name1 := "subID"
	serial1 := int64(1)
	name2, serial2, ok := subIDDecode(subIDEncode(name1, serial1))
	assert.True(this.T(), ok)
	assert.Equal(this.T(), name1, name2)
	assert.Equal(this.T(), serial1, serial2)

	_, _, ok = subIDDecode("!?")
	assert.False(this.T(), ok)
}
