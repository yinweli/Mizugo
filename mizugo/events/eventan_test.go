package events

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

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
	this.Change("test-event")
}

func (this *SuiteEventan) TearDownSuite() {
	this.Restore()
}

func (this *SuiteEventan) TestEventan() {
	awake := atomic.Bool{}
	start := atomic.Bool{}
	dispose := atomic.Bool{}
	update := atomic.Int64{}
	target := NewEventan(func(event any) {
		if e, ok := event.(*Awake); ok {
			if e.Param.(string) == "awake" {
				awake.Store(true)
			} // if
		} // if

		if e, ok := event.(*Start); ok {
			if e.Param.(string) == "start" {
				start.Store(true)
			} // if
		} // if

		if e, ok := event.(*Dispose); ok {
			if e.Param.(string) == "dispose" {
				dispose.Store(true)
			} // if
		} // if

		if e, ok := event.(*Update); ok {
			if e.Param.(string) == "update" {
				update.Add(1)
			} // if
		} // if
	})
	target.Initialize()
	target.InvokeAwake("awake")
	target.InvokeStart("start")
	target.InvokeDispose("dispose")
	target.InvokeUpdate("update", time.Millisecond*100)
	time.Sleep(time.Second)
	target.Finalize()
	time.Sleep(time.Millisecond * 100)

	assert.True(this.T(), awake.Load())
	assert.True(this.T(), start.Load())
	assert.True(this.T(), dispose.Load())
	assert.Greater(this.T(), update.Load(), int64(1))
}
