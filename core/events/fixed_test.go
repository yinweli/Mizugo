package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestFixed(t *testing.T) {
	suite.Run(t, new(SuiteFixed))
}

type SuiteFixed struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteFixed) SetupSuite() {
	this.Change("test-events-fixed")
}

func (this *SuiteFixed) TearDownSuite() {
	this.Restore()
}

func (this *SuiteFixed) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteFixed) TestFixed() {
	target := &Fixed{}
	assert.False(this.T(), target.State())
	target.Stop()
	assert.True(this.T(), target.State())
}
