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
	testdata.Env
}

func (this *SuiteEventmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-events-eventmgr")
}

func (this *SuiteEventmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteEventmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteEventmgr) TestEventmgr() {
	target := NewEventmgr(100)
	assert.NotNil(this.T(), target)
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize())
	assert.NotNil(this.T(), target.Initialize()) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
}

func (this *SuiteEventmgr) TestPubOnce() {
	target := NewEventmgr(100)
	assert.Nil(this.T(), target.Initialize())

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
	target.PubOnce(value, value)
	target.Unsub("")
}

func (this *SuiteEventmgr) TestPubDelay() {
	target := NewEventmgr(100)
	assert.Nil(this.T(), target.Initialize())

	value := "value delay"
	valid := atomic.Bool{}
	subID := ""
	subID = target.Sub(value, func(param any) {
		valid.Store(param.(string) == value)
		target.Unsub(subID)
	})

	target.PubDelay(value, value, testdata.Timeout)
	time.Sleep(testdata.Timeout * 2) // 多等一下讓延遲事件發生
	assert.True(this.T(), valid.Load())

	target.Finalize()
	target.PubDelay(value, value, testdata.Timeout)
	target.Unsub("")
}

func (this *SuiteEventmgr) TestPubFixed() {
	target := NewEventmgr(100)
	assert.Nil(this.T(), target.Initialize())

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
	target.PubFixed(value, value, testdata.Timeout)
	target.Unsub("")
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

	value1 := "value cycle1"
	value2 := "value cycle2"
	valid1 := 0
	valid2 := 0
	target.sub(value1, func(param any) {
		valid1++
		target.sub(value2, func(param any) {
			valid2++
		})
	})
	target.pub(value1, nil)
	target.pub(value2, nil)
	assert.Equal(this.T(), 1, valid1)
	assert.Equal(this.T(), 1, valid2)
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
