package helps

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestCurrent(t *testing.T) {
	suite.Run(t, new(SuiteCurrent))
}

type SuiteCurrent struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteCurrent) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-current")
}

func (this *SuiteCurrent) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteCurrent) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteCurrent) TestCurrent() {
	target := NewCurrent()
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Curr())
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
	target.SetBaseNow()
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
	target.SetBaseTime(Time())
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
	target.SetBaseDate(2023, 2, 10, 0, 0, 0)
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
	target.SetBaseDay(2023, 2, 10)
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
	target.AddBaseTime(TimeMillisecond)
	assert.WithinDuration(this.T(), target.GetBase(), target.GetTime(), testdata.Timeout)
}
