package events

import (
    `sync/atomic`
    `testing`
    `time`

    `github.com/stretchr/testify/assert`
    `github.com/stretchr/testify/suite`

    `github.com/yinweli/Mizugo/testdata`
)

func TestEvent(t *testing.T) {
    suite.Run(t, new(SuiteEvent))
}

type SuiteEvent struct {
    suite.Suite
    testdata.TestEnv
}

func (this *SuiteEvent) SetupSuite() {
    this.Change("test-event")
}

func (this *SuiteEvent) TearDownSuite() {
    this.Restore()
}

func (this *SuiteEvent) TestNewEvent() {
    assert.NotNil(this.T(), NewEvent(1))
}

func (this *SuiteEvent) TestEvent() {
    target := NewEvent(10)
    interval := time.Millisecond * 100
    awake := atomic.Bool{}
    start := atomic.Bool{}
    update := atomic.Bool{}
    dispose := atomic.Bool{}

    target.Initialize(interval, func(data Data) {
        if data.Type == Awake && data.Param == "awake" {
            awake.Store(true)
        } // if

        if data.Type == Start && data.Param == "start" {
            start.Store(true)
        } // if

        if data.Type == Update {
            update.Store(true)
        } // if

        if data.Type == Dispose && data.Param == "dispose" {
            dispose.Store(true)
        } // if
    })
    target.Execute(Awake, "awake")
    target.Execute(Start, "start")
    target.Execute(Dispose, "dispose")
    time.Sleep(interval * 2)
    target.Finalize()
    time.Sleep(interval)
    assert.True(this.T(), awake.Load())
    assert.True(this.T(), start.Load())
    assert.True(this.T(), update.Load())
    assert.True(this.T(), dispose.Load())
}
