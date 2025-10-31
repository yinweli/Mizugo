package trials

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func TestTime(t *testing.T) {
	suite.Run(t, new(SuiteTime))
}

type SuiteTime struct {
	suite.Suite
}

func (this *SuiteTime) TestWaitTimeout() {
	now := time.Now()
	WaitTimeout()
	this.GreaterOrEqual(time.Since(now), Timeout)

	now = time.Now()
	WaitTimeout(time.Second)
	this.GreaterOrEqual(time.Since(now), time.Second)
}

func (this *SuiteTime) TestWaitFor() {
	now := time.Now()
	WaitFor(time.Second, func() bool {
		return true
	})
	this.LessOrEqual(time.Since(now), time.Second)

	now = time.Now()
	WaitFor(time.Second, func() bool {
		return false
	})
	this.GreaterOrEqual(time.Since(now), time.Second)
}
