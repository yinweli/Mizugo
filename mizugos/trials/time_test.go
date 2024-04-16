package trials

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.GreaterOrEqual(this.T(), time.Since(now), Timeout)

	now = time.Now()
	WaitTimeout(time.Second)
	assert.GreaterOrEqual(this.T(), time.Since(now), time.Second)
}
