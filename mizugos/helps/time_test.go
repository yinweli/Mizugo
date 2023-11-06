package helps

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTime(t *testing.T) {
	suite.Run(t, new(SuiteTime))
}

type SuiteTime struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteTime) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-time")
}

func (this *SuiteTime) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteTime) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteTime) TestTimeZone() {
	assert.Equal(this.T(), time.UTC, GetTimeZone())
	assert.Nil(this.T(), SetTimeZone("Asia/Taipei"))
	assert.NotNil(this.T(), GetTimeZone())
	assert.NotNil(this.T(), SetTimeZone(testdata.Unknown))
}

func (this *SuiteTime) TestTime() {
	t := Time()
	assert.NotNil(this.T(), t)
	fmt.Println(t)
}

func (this *SuiteTime) TestDate() {
	t := Date(2023, 1, 1, 0, 0, 0, 0)
	assert.NotNil(this.T(), t)
	fmt.Println(t)
}

func (this *SuiteTime) TestBetween() {
	now := Time()
	start := now.Add(-TimeMinute)
	end := now.Add(TimeMinute)
	assert.True(this.T(), Between(start, end, now, true))
	assert.True(this.T(), Between(start, time.Time{}, now, true))
	assert.True(this.T(), Between(time.Time{}, end, now, true))
	assert.False(this.T(), Between(start, end, now.Add(TimeHour), true))
	assert.True(this.T(), Between(time.Time{}, time.Time{}, now, true))
	assert.False(this.T(), Between(time.Time{}, time.Time{}, now, false))
}

func (this *SuiteTime) TestBetweenf() {
	now := Time()
	start := now.Add(-TimeMinute).Format(LayoutSecond)
	end := now.Add(TimeMinute).Format(LayoutSecond)
	assert.True(this.T(), Betweenf(LayoutSecond, start, end, now, true))
	assert.True(this.T(), Betweenf(LayoutSecond, start, "", now, true))
	assert.True(this.T(), Betweenf(LayoutSecond, "", end, now, true))
	assert.False(this.T(), Betweenf(LayoutSecond, start, end, now.Add(TimeHour), true))
	assert.True(this.T(), Betweenf(LayoutSecond, "", "", now, true))
	assert.False(this.T(), Betweenf(LayoutSecond, "", "", now, false))
	assert.False(this.T(), Betweenf(testdata.Unknown, start, end, now, true))
	assert.False(this.T(), Betweenf(LayoutSecond, start, testdata.Unknown, now.Add(TimeHour), true))
	assert.False(this.T(), Betweenf(LayoutSecond, testdata.Unknown, end, now.Add(TimeHour), true))
}

func (this *SuiteTime) TestSameDay() {
	assert.True(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2023, 6, 1, 0, 0, 0, 0)))
	assert.True(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2023, 6, 1, 1, 0, 0, 0)))
	assert.True(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2023, 6, 1, 0, 1, 0, 0)))
	assert.True(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2023, 6, 1, 0, 0, 1, 0)))
	assert.False(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2023, 6, 2, 0, 0, 0, 0)))
	assert.False(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2024, 7, 1, 0, 0, 0, 0)))
	assert.False(this.T(), SameDay(
		Date(2023, 6, 1, 0, 0, 0, 0),
		Date(2024, 6, 1, 0, 0, 0, 0)))
}

func (this *SuiteTime) TestDaily() {
	assert.True(this.T(), Daily(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 1, 1, 12, 0, 0, 0), 12))
	assert.False(this.T(), Daily(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 3, 1, 12, 0, 0, 0), 12))
	assert.False(this.T(), Daily(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 12, 0, 0, 0), 12))
	assert.True(this.T(), Daily(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 11, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		DailyPrev(Date(2023, 2, 1, 12, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 1, 31, 12, 0, 0, 0),
		DailyPrev(Date(2023, 2, 1, 11, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		DailyPrev(Date(2023, 2, 1, 13, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 2, 12, 0, 0, 0),
		DailyNext(Date(2023, 2, 1, 12, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		DailyNext(Date(2023, 2, 1, 11, 0, 0, 0), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 2, 12, 0, 0, 0),
		DailyNext(Date(2023, 2, 1, 13, 0, 0, 0), 12))
}

func (this *SuiteTime) TestWeekly() {
	assert.True(this.T(), Weekly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 1, 1, 12, 0, 0, 0), 3, 12))
	assert.False(this.T(), Weekly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 3, 1, 12, 0, 0, 0), 3, 12))
	assert.False(this.T(), Weekly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 12, 0, 0, 0), 3, 12))
	assert.True(this.T(), Weekly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 11, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		WeeklyPrev(Date(2023, 2, 1, 12, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 25, 12, 0, 0, 0),
		WeeklyPrev(Date(2023, 2, 1, 11, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		WeeklyPrev(Date(2023, 2, 1, 13, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		WeeklyPrev(Date(2023, 2, 2, 12, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 25, 12, 0, 0, 0),
		WeeklyPrev(Date(2023, 1, 31, 12, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12, 0, 0, 0),
		WeeklyNext(Date(2023, 2, 1, 12, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		WeeklyNext(Date(2023, 2, 1, 11, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12, 0, 0, 0),
		WeeklyNext(Date(2023, 2, 1, 13, 0, 0, 0), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12, 0, 0, 0),
		WeeklyNext(Date(2023, 2, 2, 12, 0, 0, 0), 3, 12))
}

func (this *SuiteTime) TestMonthly() {
	assert.True(this.T(), Monthly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 1, 1, 12, 0, 0, 0), 1, 12))
	assert.False(this.T(), Monthly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 3, 1, 12, 0, 0, 0), 1, 12))
	assert.False(this.T(), Monthly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 12, 0, 0, 0), 1, 12))
	assert.True(this.T(), Monthly(
		Date(2023, 2, 1, 12, 0, 0, 0),
		Date(2023, 2, 1, 11, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		MonthlyPrev(Date(2023, 2, 1, 12, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 12, 0, 0, 0),
		MonthlyPrev(Date(2023, 2, 1, 11, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		MonthlyPrev(Date(2023, 2, 1, 13, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 28, 12, 0, 0, 0),
		MonthlyPrev(Date(2023, 3, 1, 13, 0, 0, 0), 31, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 31, 12, 0, 0, 0),
		MonthlyPrev(Date(2023, 4, 1, 13, 0, 0, 0), 31, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 1, 12, 0, 0, 0),
		MonthlyNext(Date(2023, 2, 1, 12, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12, 0, 0, 0),
		MonthlyNext(Date(2023, 2, 1, 11, 0, 0, 0), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 1, 12, 0, 0, 0),
		MonthlyNext(Date(2023, 2, 1, 13, 0, 0, 0), 1, 12))
}

func (this *SuiteTime) TestFixedPrev() {
	duration := TimeHour
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 23, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 23, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 23, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 50, 0, 0), duration))
	duration = TimeHour * 2
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 22, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 22, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 22, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 50, 0, 0), duration))
	duration = TimeHour * 3
	assert.Equal(this.T(),
		Date(2023, 1, 1, 6, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 6, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 6, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 21, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 21, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 21, 0, 0, 0),
		FixedPrev(Date(2023, 1, 1, 23, 50, 0, 0), duration))
}

func (this *SuiteTime) TestFixedNext() {
	duration := TimeHour
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 50, 0, 0), duration))
	duration = TimeHour * 2
	assert.Equal(this.T(),
		Date(2023, 1, 1, 10, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 10, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 10, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 50, 0, 0), duration))
	duration = TimeHour * 3
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 8, 50, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 0, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 10, 0, 0), duration))
	assert.Equal(this.T(),
		Date(2023, 1, 2, 0, 0, 0, 0),
		FixedNext(Date(2023, 1, 1, 23, 50, 0, 0), duration))
}

func (this *SuiteTime) TestFixedCheck() {
	assert.True(this.T(), FixedCheck(TimeHour))
	assert.True(this.T(), FixedCheck(TimeHour*2))
	assert.True(this.T(), FixedCheck(TimeHour*3))
	assert.True(this.T(), FixedCheck(TimeHour*4))
	assert.False(this.T(), FixedCheck(TimeHour*5))
	assert.True(this.T(), FixedCheck(TimeHour*6))
	assert.False(this.T(), FixedCheck(TimeHour*7))
	assert.True(this.T(), FixedCheck(TimeHour*8))
	assert.False(this.T(), FixedCheck(TimeHour*9))
}
