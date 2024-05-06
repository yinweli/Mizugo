package helps

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestTime(t *testing.T) {
	suite.Run(t, new(SuiteTime))
}

type SuiteTime struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTime) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-time"))
}

func (this *SuiteTime) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTime) TestTimeZone() {
	assert.Equal(this.T(), time.UTC, GetTimeZone())
	SetTimeZoneUTC()
	assert.Equal(this.T(), time.UTC, GetTimeZone())
	SetTimeZoneLocal()
	assert.Equal(this.T(), time.Local, GetTimeZone())
	assert.Nil(this.T(), SetTimeZone("Asia/Taipei"))
	assert.NotNil(this.T(), GetTimeZone())
	assert.NotNil(this.T(), SetTimeZone(testdata.Unknown))
}

func (this *SuiteTime) TestTime() {
	now := Time()
	assert.NotNil(this.T(), now)
	fmt.Println(now)
}

func (this *SuiteTime) TestTimef() {
	now, err := Timef(LayoutSecond, "2023-02-05 01:02:03")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now, err = Timef(LayoutMinute, "2023-02-05 01:02")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now, err = Timef(LayoutDay, "2023-02-05")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now, err = Timef(LayoutDay, "")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	_, err = Timef(testdata.Unknown, "2023-02-05")
	assert.NotNil(this.T(), err)
}

func (this *SuiteTime) TestDate() {
	now := Date(2023, 2, 5, 1, 2, 3, 4)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1, 2, 3)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1, 2)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
	now = Date(2023, 2, 5)
	assert.NotNil(this.T(), now)
	fmt.Println(now)
}

func (this *SuiteTime) TestBefore() {
	assert.True(this.T(), Before(Date(2023, 2, 1), Date(2023, 2, 5)))
	assert.False(this.T(), Before(Date(2023, 2, 10), Date(2023, 2, 5)))

	assert.True(this.T(), Beforef(LayoutDay, "2023-02-01", Date(2023, 2, 5)))
	assert.False(this.T(), Beforef(LayoutDay, "2023-02-10", Date(2023, 2, 5)))
	assert.False(this.T(), Beforef(LayoutDay, testdata.Unknown, Date(2023, 2, 5)))

	assert.True(this.T(), Beforefx(LayoutDay, "2023-02-01", "2023-02-05"))
	assert.False(this.T(), Beforefx(LayoutDay, "2023-02-10", "2023-02-05"))
	assert.False(this.T(), Beforefx(LayoutDay, testdata.Unknown, "2023-02-05"))
	assert.False(this.T(), Beforefx(LayoutDay, "2023-02-10", testdata.Unknown))
}

func (this *SuiteTime) TestAfter() {
	assert.True(this.T(), After(Date(2023, 2, 10), Date(2023, 2, 5)))
	assert.False(this.T(), After(Date(2023, 2, 1), Date(2023, 2, 5)))

	assert.True(this.T(), Afterf(LayoutDay, "2023-02-10", Date(2023, 2, 5)))
	assert.False(this.T(), Afterf(LayoutDay, "2023-02-1", Date(2023, 2, 5)))
	assert.False(this.T(), Afterf(LayoutDay, testdata.Unknown, Date(2023, 2, 5)))

	assert.True(this.T(), Afterfx(LayoutDay, "2023-02-10", "2023-02-05"))
	assert.False(this.T(), Afterfx(LayoutDay, "2023-02-01", "2023-02-05"))
	assert.False(this.T(), Afterfx(LayoutDay, testdata.Unknown, "2023-02-05"))
	assert.False(this.T(), Afterfx(LayoutDay, "2023-02-10", testdata.Unknown))
}

func (this *SuiteTime) TestBetween() {
	assert.True(this.T(), Between(Date(2023, 2, 5), Date(2023, 2, 15), Date(2023, 2, 10)))
	assert.False(this.T(), Between(Date(2023, 2, 5), Date(2023, 2, 15), Date(2023, 2, 20)))
	assert.True(this.T(), Between(time.Time{}, time.Time{}, Date(2023, 2, 10)))
	assert.False(this.T(), Between(time.Time{}, time.Time{}, Date(2023, 2, 10), false))
	assert.True(this.T(), Between(time.Time{}, Date(2023, 2, 15), Date(2023, 2, 10)))
	assert.False(this.T(), Between(time.Time{}, Date(2023, 2, 15), Date(2023, 2, 20)))
	assert.True(this.T(), Between(Date(2023, 2, 5), time.Time{}, Date(2023, 2, 10)))
	assert.False(this.T(), Between(Date(2023, 2, 5), time.Time{}, Date(2023, 2, 1)))

	assert.True(this.T(), Betweenf(LayoutDay, "2023-02-05", "2023-02-15", Date(2023, 2, 10)))
	assert.False(this.T(), Betweenf(LayoutDay, "2023-02-05", "2023-02-15", Date(2023, 2, 20)))
	assert.False(this.T(), Betweenf(LayoutDay, testdata.Unknown, "2023-02-15", Date(2023, 2, 10)))
	assert.False(this.T(), Betweenf(LayoutDay, "2023-02-05", testdata.Unknown, Date(2023, 2, 10)))
	assert.True(this.T(), Betweenf(LayoutDay, "", "", Date(2023, 2, 10)))
	assert.False(this.T(), Betweenf(LayoutDay, "", "", Date(2023, 2, 10), false))
	assert.True(this.T(), Betweenf(LayoutDay, "", "2023-02-15", Date(2023, 2, 10)))
	assert.False(this.T(), Betweenf(LayoutDay, "", "2023-02-15", Date(2023, 2, 20)))
	assert.True(this.T(), Betweenf(LayoutDay, "2023-02-05", "", Date(2023, 2, 10)))
	assert.False(this.T(), Betweenf(LayoutDay, "2023-02-05", "", Date(2023, 2, 1)))

	assert.True(this.T(), Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", "2023-02-20"))
	assert.False(this.T(), Betweenfx(LayoutDay, testdata.Unknown, "2023-02-15", "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "2023-02-05", testdata.Unknown, "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", testdata.Unknown))
	assert.True(this.T(), Betweenfx(LayoutDay, "", "", "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "", "", "2023-02-10", false))
	assert.True(this.T(), Betweenfx(LayoutDay, "", "2023-02-15", "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "", "2023-02-15", "2023-02-20"))
	assert.True(this.T(), Betweenfx(LayoutDay, "2023-02-05", "", "2023-02-10"))
	assert.False(this.T(), Betweenfx(LayoutDay, "2023-02-05", "", "2023-02-01"))
}

func (this *SuiteTime) TestOverlap() {
	assert.True(this.T(), Overlap(
		Date(2023, 2, 1), Date(2023, 2, 10),
		Date(2023, 2, 5), Date(2023, 2, 15)))
	assert.False(this.T(), Overlap(
		Date(2023, 2, 1), Date(2023, 2, 10),
		Date(2023, 2, 15), Date(2023, 2, 20)))

	assert.True(this.T(), Overlapf(LayoutDay,
		"2023-02-01", "2023-02-10",
		Date(2023, 2, 5), Date(2023, 2, 15)))
	assert.False(this.T(), Overlapf(LayoutDay,
		"2023-02-01", "2023-02-10",
		Date(2023, 2, 15), Date(2023, 2, 20)))
	assert.False(this.T(), Overlapf(LayoutDay,
		testdata.Unknown, "2023-02-10",
		Date(2023, 2, 15), Date(2023, 2, 20)))
	assert.False(this.T(), Overlapf(LayoutDay,
		"2023-02-01", testdata.Unknown,
		Date(2023, 2, 15), Date(2023, 2, 20)))

	assert.True(this.T(), Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-05", "2023-02-15"))
	assert.False(this.T(), Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-15", "2023-02-20"))
	assert.False(this.T(), Overlapfx(LayoutDay,
		testdata.Unknown, "2023-02-10",
		"2023-02-15", "2023-02-20"))
	assert.False(this.T(), Overlapfx(LayoutDay,
		"2023-02-01", testdata.Unknown,
		"2023-02-15", "2023-02-20"))
	assert.False(this.T(), Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		testdata.Unknown, "2023-02-20"))
	assert.False(this.T(), Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-15", testdata.Unknown))
}

func (this *SuiteTime) TestDaily() {
	assert.True(this.T(), Daily(
		Date(2023, 2, 1, 12),
		Date(2023, 1, 1, 12), 12))
	assert.False(this.T(), Daily(
		Date(2023, 2, 1, 12),
		Date(2023, 3, 1, 12), 12))
	assert.False(this.T(), Daily(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 12), 12))
	assert.True(this.T(), Daily(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 11), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		DailyPrev(Date(2023, 2, 1, 12), 12))
	assert.Equal(this.T(),
		Date(2023, 1, 31, 12),
		DailyPrev(Date(2023, 2, 1, 11), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		DailyPrev(Date(2023, 2, 1, 13), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 2, 12),
		DailyNext(Date(2023, 2, 1, 12), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		DailyNext(Date(2023, 2, 1, 11), 12))
	assert.Equal(this.T(),
		Date(2023, 2, 2, 12),
		DailyNext(Date(2023, 2, 1, 13), 12))
}

func (this *SuiteTime) TestWeekly() {
	assert.True(this.T(), Weekly(
		Date(2023, 2, 1, 12),
		Date(2023, 1, 1, 12), 3, 12))
	assert.False(this.T(), Weekly(
		Date(2023, 2, 1, 12),
		Date(2023, 3, 1, 12), 3, 12))
	assert.False(this.T(), Weekly(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 12), 3, 12))
	assert.True(this.T(), Weekly(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 11), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		WeeklyPrev(Date(2023, 2, 1, 12), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 25, 12),
		WeeklyPrev(Date(2023, 2, 1, 11), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		WeeklyPrev(Date(2023, 2, 1, 13), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		WeeklyPrev(Date(2023, 2, 2, 12), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 25, 12),
		WeeklyPrev(Date(2023, 1, 31, 12), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12),
		WeeklyNext(Date(2023, 2, 1, 12), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		WeeklyNext(Date(2023, 2, 1, 11), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12),
		WeeklyNext(Date(2023, 2, 1, 13), 3, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 8, 12),
		WeeklyNext(Date(2023, 2, 2, 12), 3, 12))
}

func (this *SuiteTime) TestMonthly() {
	assert.True(this.T(), Monthly(
		Date(2023, 2, 1, 12),
		Date(2023, 1, 1, 12), 1, 12))
	assert.False(this.T(), Monthly(
		Date(2023, 2, 1, 12),
		Date(2023, 3, 1, 12), 1, 12))
	assert.False(this.T(), Monthly(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 12), 1, 12))
	assert.True(this.T(), Monthly(
		Date(2023, 2, 1, 12),
		Date(2023, 2, 1, 11), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		MonthlyPrev(Date(2023, 2, 1, 12), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 12),
		MonthlyPrev(Date(2023, 2, 1, 11), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		MonthlyPrev(Date(2023, 2, 1, 13), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 28, 12),
		MonthlyPrev(Date(2023, 3, 1, 13), 31, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 31, 12),
		MonthlyPrev(Date(2023, 4, 1, 13), 31, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 1, 12),
		MonthlyNext(Date(2023, 2, 1, 12), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 2, 1, 12),
		MonthlyNext(Date(2023, 2, 1, 11), 1, 12))
	assert.Equal(this.T(),
		Date(2023, 3, 1, 12),
		MonthlyNext(Date(2023, 2, 1, 13), 1, 12))
}

func (this *SuiteTime) TestFixedPrev() {
	assert.Equal(this.T(),
		Date(1970, 1, 1),
		FixedPrev(time.Time{}, time.Time{}, TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8),
		FixedPrev(time.Time{}, Date(2023, 1, 1, 8), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8),
		FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8),
		FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 10), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8),
		FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 50), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 7),
		FixedPrev(Date(1970, 1, 1, 1), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 8),
		FixedPrev(Date(1970, 1, 1, 2), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 6, 30),
		FixedPrev(Date(1970, 1, 1, 3, 30), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
}

func (this *SuiteTime) TestFixedNext() {
	assert.Equal(this.T(),
		Date(1970, 1, 1, 1),
		FixedNext(time.Time{}, time.Time{}, TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9),
		FixedNext(time.Time{}, Date(2023, 1, 1, 8), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9),
		FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9),
		FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 10), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9),
		FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 50), TimeHour))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 10),
		FixedNext(Date(1970, 1, 1, 1), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 11),
		FixedNext(Date(1970, 1, 1, 2), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	assert.Equal(this.T(),
		Date(2023, 1, 1, 9, 30),
		FixedNext(Date(1970, 1, 1, 3, 30), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
}
