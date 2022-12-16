package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestWaitTimeout(t *testing.T) {
	suite.Run(t, new(SuiteWaitTimeout))
}

type SuiteWaitTimeout struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteWaitTimeout) SetupSuite() {
	this.Change("test-utils-waitTimeout")
}

func (this *SuiteWaitTimeout) TearDownSuite() {
	this.Restore()
}

func (this *SuiteWaitTimeout) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteWaitTimeout) TestNewWaitTimeout() {
	assert.NotNil(this.T(), NewWaitTimeout(0))
}

func (this *SuiteWaitTimeout) TestWaitTimeout() {
	target := NewWaitTimeout(testdata.Timeout)
	target.Done()
	assert.True(this.T(), target.Wait())

	startTime := time.Now()
	target = NewWaitTimeout(testdata.Timeout)
	assert.False(this.T(), target.Wait())
	assert.GreaterOrEqual(this.T(), time.Now().Add(-testdata.Timeout), startTime)
}
