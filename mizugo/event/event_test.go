package event

import (
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
    valid := false

    target.Begin(func(event any) {
        if s, ok := event.(string); ok {
            if s == "event" {
                valid = true
            } // if
        } // if
    })
    target.Add("event")
    time.Sleep(time.Millisecond * 10)
    target.End()
    time.Sleep(time.Millisecond * 10)
    assert.True(this.T(), valid)
}
