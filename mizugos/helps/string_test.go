package helps

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestString(t *testing.T) {
	suite.Run(t, new(SuiteString))
}

type SuiteString struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteString) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-string"))
}

func (this *SuiteString) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteString) TestStringDisplayLength() {
	this.Equal(30, StringDisplayLength("Hello, こんにちは, 안녕하세요!"))
}

func (this *SuiteString) TestStringPercentage() {
	this.Equal("1.00%", StringPercentage(1, 100))
	this.Equal("0.00%", StringPercentage(1, 0))
}

func (this *SuiteString) TestStringDuration() {
	target, err := StringDuration("1d 1h 1m 1s 1ms")
	this.Nil(err)
	this.Equal(TimeDay+TimeHour+TimeMinute+TimeSecond+TimeMillisecond, target)
	target, err = StringDuration("- 1d 1h 1m 1s 1ms")
	this.Nil(err)
	this.Equal(-(TimeDay + TimeHour + TimeMinute + TimeSecond + TimeMillisecond), target)
	target, err = StringDuration("  ")
	this.Nil(err)
	this.Equal(time.Duration(0), target)
	_, err = StringDuration("?d ?h ?m ?s ?ms")
	this.NotNil(err)
	_, err = StringDuration("1? 2? 3? 4? 5??")
	this.NotNil(err)
	_, err = StringDuration("????????????????????")
	this.NotNil(err)
}
