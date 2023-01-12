package events

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestEventmgr(t *testing.T) {
	suite.Run(t, new(SuiteEventmgr))
}

type SuiteEventmgr struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteEventmgr) SetupSuite() {
	this.Change("test-events-eventmgr")
}

func (this *SuiteEventmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEventmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
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
	value := "value once"
	valid := atomic.Bool{}
	subID := ""
	subID = target.Sub(value, func(param any) {
		valid.Store(param.(string) == value)
		target.Unsub(subID)
	})

	target.PubOnce(value, value)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubOnce(value, value) // 測試在結束之後發布事件
	target.Unsub("")             // 測試在結束之後取消訂閱
}

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(10)
	target.Initialize()
	value := "value fixed"
	valid := atomic.Bool{}
	subID := ""
	subID = target.Sub(value, func(param any) {
		valid.Store(param.(string) == value)
		target.Unsub(subID)
	})

	target.PubFixed(value, value, testdata.Timeout)
	time.Sleep(testdata.Timeout * 2) // 多等一下讓定時事件發生
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubFixed(value, value, testdata.Timeout) // 測試在結束之後發布事件
	target.Unsub("")                                // 測試在結束之後取消訂閱
}

func (this *SuiteEventmgr) TestPubsub() {
	target := newPubsub()
	value := "value pubsub"
	valid := 0
	validFunc := func(param any) {
		if param.(string) == value {
			valid++
		} // if
	}
	subID1 := target.sub(value, validFunc)
	subID2 := target.sub(value, validFunc)

	target.pub(value, value)
	assert.Equal(this.T(), 2, valid)

	target.unsub(subID1)
	target.pub(value, value)
	assert.Equal(this.T(), 3, valid)

	target.unsub(subID2)
	target.pub(value, value)
	assert.Equal(this.T(), 3, valid)
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
