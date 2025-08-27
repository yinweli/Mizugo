package helps

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.Equal(time.UTC, GetTimeZone())
	SetTimeZoneUTC()
	this.Equal(time.UTC, GetTimeZone())
	SetTimeZoneLocal()
	this.Equal(time.Local, GetTimeZone())
	this.Nil(SetTimeZone("Asia/Taipei"))
	this.NotNil(GetTimeZone())
	this.NotNil(SetTimeZone(testdata.Unknown))
}

func (this *SuiteTime) TestTime() {
	now := Time()
	this.NotNil(now)
	fmt.Println(now)
}

func (this *SuiteTime) TestTimef() {
	now, err := Timef(LayoutSecond, "2023-02-05 01:02:03")
	this.Nil(err)
	this.NotNil(now)
	fmt.Println(now)
	now, err = Timef(LayoutMinute, "2023-02-05 01:02")
	this.Nil(err)
	this.NotNil(now)
	fmt.Println(now)
	now, err = Timef(LayoutDay, "2023-02-05")
	this.Nil(err)
	this.NotNil(now)
	fmt.Println(now)
	now, err = Timef(LayoutDay, "")
	this.Nil(err)
	this.NotNil(now)
	fmt.Println(now)
	_, err = Timef(testdata.Unknown, "2023-02-05")
	this.NotNil(err)
}

func (this *SuiteTime) TestDate() {
	now := Date(2023, 2, 5, 1, 2, 3, 4)
	this.NotNil(now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1, 2, 3)
	this.NotNil(now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1, 2)
	this.NotNil(now)
	fmt.Println(now)
	now = Date(2023, 2, 5, 1)
	this.NotNil(now)
	fmt.Println(now)
	now = Date(2023, 2, 5)
	this.NotNil(now)
	fmt.Println(now)
}

func (this *SuiteTime) TestBefore() {
	this.True(Before(Date(2023, 2, 1), Date(2023, 2, 5)))
	this.False(Before(Date(2023, 2, 10), Date(2023, 2, 5)))

	this.True(Beforef(LayoutDay, "2023-02-01", Date(2023, 2, 5)))
	this.False(Beforef(LayoutDay, "2023-02-10", Date(2023, 2, 5)))
	this.False(Beforef(LayoutDay, testdata.Unknown, Date(2023, 2, 5)))

	this.True(Beforefx(LayoutDay, "2023-02-01", "2023-02-05"))
	this.False(Beforefx(LayoutDay, "2023-02-10", "2023-02-05"))
	this.False(Beforefx(LayoutDay, testdata.Unknown, "2023-02-05"))
	this.False(Beforefx(LayoutDay, "2023-02-10", testdata.Unknown))
}

func (this *SuiteTime) TestAfter() {
	this.True(After(Date(2023, 2, 10), Date(2023, 2, 5)))
	this.False(After(Date(2023, 2, 1), Date(2023, 2, 5)))

	this.True(Afterf(LayoutDay, "2023-02-10", Date(2023, 2, 5)))
	this.False(Afterf(LayoutDay, "2023-02-1", Date(2023, 2, 5)))
	this.False(Afterf(LayoutDay, testdata.Unknown, Date(2023, 2, 5)))

	this.True(Afterfx(LayoutDay, "2023-02-10", "2023-02-05"))
	this.False(Afterfx(LayoutDay, "2023-02-01", "2023-02-05"))
	this.False(Afterfx(LayoutDay, testdata.Unknown, "2023-02-05"))
	this.False(Afterfx(LayoutDay, "2023-02-10", testdata.Unknown))
}

func (this *SuiteTime) TestBetween() {
	this.True(Between(Date(2023, 2, 5), Date(2023, 2, 15), Date(2023, 2, 10)))
	this.False(Between(Date(2023, 2, 5), Date(2023, 2, 15), Date(2023, 2, 20)))
	this.True(Between(time.Time{}, time.Time{}, Date(2023, 2, 10)))
	this.False(Between(time.Time{}, time.Time{}, Date(2023, 2, 10), false))
	this.True(Between(time.Time{}, Date(2023, 2, 15), Date(2023, 2, 10)))
	this.False(Between(time.Time{}, Date(2023, 2, 15), Date(2023, 2, 20)))
	this.True(Between(Date(2023, 2, 5), time.Time{}, Date(2023, 2, 10)))
	this.False(Between(Date(2023, 2, 5), time.Time{}, Date(2023, 2, 1)))

	this.True(Betweenf(LayoutDay, "2023-02-05", "2023-02-15", Date(2023, 2, 10)))
	this.False(Betweenf(LayoutDay, "2023-02-05", "2023-02-15", Date(2023, 2, 20)))
	this.False(Betweenf(LayoutDay, testdata.Unknown, "2023-02-15", Date(2023, 2, 10)))
	this.False(Betweenf(LayoutDay, "2023-02-05", testdata.Unknown, Date(2023, 2, 10)))
	this.True(Betweenf(LayoutDay, "", "", Date(2023, 2, 10)))
	this.False(Betweenf(LayoutDay, "", "", Date(2023, 2, 10), false))
	this.True(Betweenf(LayoutDay, "", "2023-02-15", Date(2023, 2, 10)))
	this.False(Betweenf(LayoutDay, "", "2023-02-15", Date(2023, 2, 20)))
	this.True(Betweenf(LayoutDay, "2023-02-05", "", Date(2023, 2, 10)))
	this.False(Betweenf(LayoutDay, "2023-02-05", "", Date(2023, 2, 1)))

	this.True(Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", "2023-02-20"))
	this.False(Betweenfx(LayoutDay, testdata.Unknown, "2023-02-15", "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "2023-02-05", testdata.Unknown, "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "2023-02-05", "2023-02-15", testdata.Unknown))
	this.True(Betweenfx(LayoutDay, "", "", "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "", "", "2023-02-10", false))
	this.True(Betweenfx(LayoutDay, "", "2023-02-15", "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "", "2023-02-15", "2023-02-20"))
	this.True(Betweenfx(LayoutDay, "2023-02-05", "", "2023-02-10"))
	this.False(Betweenfx(LayoutDay, "2023-02-05", "", "2023-02-01"))
}

func (this *SuiteTime) TestOverlap() {
	this.True(Overlap(
		Date(2023, 2, 1), Date(2023, 2, 10),
		Date(2023, 2, 5), Date(2023, 2, 15),
	))
	this.False(Overlap(
		Date(2023, 2, 1), Date(2023, 2, 10),
		Date(2023, 2, 15), Date(2023, 2, 20),
	))

	this.True(Overlapf(LayoutDay,
		"2023-02-01", "2023-02-10",
		Date(2023, 2, 5), Date(2023, 2, 15),
	))
	this.False(Overlapf(LayoutDay,
		"2023-02-01", "2023-02-10",
		Date(2023, 2, 15), Date(2023, 2, 20),
	))
	this.False(Overlapf(LayoutDay,
		testdata.Unknown, "2023-02-10",
		Date(2023, 2, 15), Date(2023, 2, 20),
	))
	this.False(Overlapf(LayoutDay,
		"2023-02-01", testdata.Unknown,
		Date(2023, 2, 15), Date(2023, 2, 20),
	))

	this.True(Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-05", "2023-02-15",
	))
	this.False(Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-15", "2023-02-20",
	))
	this.False(Overlapfx(LayoutDay,
		testdata.Unknown, "2023-02-10",
		"2023-02-15", "2023-02-20",
	))
	this.False(Overlapfx(LayoutDay,
		"2023-02-01", testdata.Unknown,
		"2023-02-15", "2023-02-20",
	))
	this.False(Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		testdata.Unknown, "2023-02-20",
	))
	this.False(Overlapfx(LayoutDay,
		"2023-02-01", "2023-02-10",
		"2023-02-15", testdata.Unknown,
	))
}

func (this *SuiteTime) TestDaily() {
	this.True(Daily(Date(2023, 2, 1, 12), Date(2023, 1, 1, 12), 12))
	this.True(Daily(Date(2023, 2, 1, 12), Date(2023, 2, 1, 11), 12))
	this.False(Daily(Date(2023, 2, 1, 12), Date(2023, 2, 1, 12), 12))
	this.False(Daily(Date(2023, 2, 1, 12), Date(2023, 3, 1, 12), 12))
}

func (this *SuiteTime) TestDailyPrev() {
	this.Equal(Date(2023, 1, 31, 12), DailyPrev(Date(2023, 2, 1, 11), 12))
	this.Equal(Date(2023, 2, 1, 12), DailyPrev(Date(2023, 2, 1, 12), 12))
	this.Equal(Date(2023, 2, 1, 12), DailyPrev(Date(2023, 2, 1, 13), 12))
	this.Equal(Date(2023, 2, 1, 0), DailyPrev(Date(2023, 2, 1, 1), -1))
	this.Equal(Date(2023, 2, 1, 23), DailyPrev(Date(2023, 2, 2, 0), 99))
}

func (this *SuiteTime) TestDailyNext() {
	this.Equal(Date(2023, 2, 1, 12), DailyNext(Date(2023, 2, 1, 11), 12))
	this.Equal(Date(2023, 2, 2, 12), DailyNext(Date(2023, 2, 1, 12), 12))
	this.Equal(Date(2023, 2, 2, 12), DailyNext(Date(2023, 2, 1, 13), 12))
	this.Equal(Date(2023, 2, 1, 0), DailyNext(Date(2023, 1, 31, 23), -1))
	this.Equal(Date(2023, 2, 1, 23), DailyNext(Date(2023, 2, 1, 22), 99))
}

func (this *SuiteTime) TestWeekly() {
	this.True(Weekly(Date(2023, 2, 1, 12), Date(2023, 1, 1, 12), 3, 12))
	this.True(Weekly(Date(2023, 2, 1, 12), Date(2023, 2, 1, 11), 3, 12))
	this.False(Weekly(Date(2023, 2, 1, 12), Date(2023, 2, 1, 12), 3, 12))
	this.False(Weekly(Date(2023, 2, 1, 12), Date(2023, 3, 1, 12), 3, 12))
}

func (this *SuiteTime) TestWeeklyPrev() {
	this.Equal(Date(2023, 1, 25, 12), WeeklyPrev(Date(2023, 1, 31, 12), 3, 12))
	this.Equal(Date(2023, 1, 25, 12), WeeklyPrev(Date(2023, 2, 1, 11), 3, 12))
	this.Equal(Date(2023, 2, 1, 12), WeeklyPrev(Date(2023, 2, 1, 12), 3, 12))
	this.Equal(Date(2023, 2, 1, 12), WeeklyPrev(Date(2023, 2, 1, 13), 3, 12))
	this.Equal(Date(2023, 2, 1, 12), WeeklyPrev(Date(2023, 2, 2, 12), 3, 12))
	this.Equal(Date(2023, 1, 25, 12), WeeklyPrev(Date(2023, 1, 31, 12), -4, 12))
	this.Equal(Date(2023, 1, 25, 12), WeeklyPrev(Date(2023, 1, 31, 12), 10, 12))
	this.Equal(Date(2023, 2, 1, 0), WeeklyPrev(Date(2023, 2, 1, 1), 3, -1))
	this.Equal(Date(2023, 2, 1, 23), WeeklyPrev(Date(2023, 2, 2, 0), 3, 99))
}

func (this *SuiteTime) TestWeeklyNext() {
	this.Equal(Date(2023, 2, 1, 12), WeeklyNext(Date(2023, 2, 1, 11), 3, 12))
	this.Equal(Date(2023, 2, 8, 12), WeeklyNext(Date(2023, 2, 1, 12), 3, 12))
	this.Equal(Date(2023, 2, 8, 12), WeeklyNext(Date(2023, 2, 1, 13), 3, 12))
	this.Equal(Date(2023, 2, 8, 12), WeeklyNext(Date(2023, 2, 2, 12), 3, 12))
	this.Equal(Date(2023, 2, 1, 12), WeeklyNext(Date(2023, 2, 1, 11), -4, 12))
	this.Equal(Date(2023, 2, 1, 12), WeeklyNext(Date(2023, 2, 1, 11), 10, 12))
	this.Equal(Date(2023, 2, 1, 0), WeeklyNext(Date(2023, 1, 31, 23), 3, -1))
	this.Equal(Date(2023, 2, 1, 23), WeeklyNext(Date(2023, 2, 1, 22), 3, 99))
}

func (this *SuiteTime) TestMonthly() {
	this.True(Monthly(Date(2023, 2, 1, 12), Date(2023, 1, 1, 12), 1, 12))
	this.True(Monthly(Date(2023, 2, 1, 12), Date(2023, 2, 1, 11), 1, 12))
	this.False(Monthly(Date(2023, 2, 1, 12), Date(2023, 2, 1, 12), 1, 12))
	this.False(Monthly(Date(2023, 2, 1, 12), Date(2023, 3, 1, 12), 1, 12))
}

func (this *SuiteTime) TestMonthlyPrev() {
	this.Equal(Date(2023, 1, 1, 12), MonthlyPrev(Date(2023, 2, 1, 11), 1, 12))
	this.Equal(Date(2023, 2, 1, 12), MonthlyPrev(Date(2023, 2, 1, 12), 1, 12))
	this.Equal(Date(2023, 2, 1, 12), MonthlyPrev(Date(2023, 2, 1, 13), 1, 12))
	this.Equal(Date(2023, 1, 31, 12), MonthlyPrev(Date(2023, 3, 1, 13), 31, 12))
	this.Equal(Date(2023, 3, 31, 12), MonthlyPrev(Date(2023, 4, 1, 13), 31, 12))
	this.Equal(Date(2023, 2, 1, 12), MonthlyPrev(Date(2023, 2, 1, 13), -1, 12))
	this.Equal(Date(2023, 1, 31, 12), MonthlyPrev(Date(2023, 2, 1, 13), 99, 12))
	this.Equal(Date(2023, 2, 1, 0), MonthlyPrev(Date(2023, 2, 1, 1), 1, -1))
	this.Equal(Date(2023, 2, 1, 23), MonthlyPrev(Date(2023, 2, 2, 0), 1, 99))
}

func (this *SuiteTime) TestMonthlyNext() {
	this.Equal(Date(2023, 2, 1, 12), MonthlyNext(Date(2023, 2, 1, 11), 1, 12))
	this.Equal(Date(2023, 3, 1, 12), MonthlyNext(Date(2023, 2, 1, 12), 1, 12))
	this.Equal(Date(2023, 3, 1, 12), MonthlyNext(Date(2023, 2, 1, 13), 1, 12))
	this.Equal(Date(2023, 3, 1, 12), MonthlyNext(Date(2023, 2, 1, 13), -1, 12))
	this.Equal(Date(2023, 3, 31, 12), MonthlyNext(Date(2023, 2, 1, 13), 99, 12))
	this.Equal(Date(2023, 2, 1, 0), MonthlyNext(Date(2023, 1, 31, 23), 1, -1))
	this.Equal(Date(2023, 2, 1, 23), MonthlyNext(Date(2023, 2, 1, 22), 1, 99))
}

func (this *SuiteTime) TestYearly() {
	this.True(Yearly(Date(2024, 2, 1, 12), Date(2023, 1, 1, 12), 2, 1, 12))
	this.True(Yearly(Date(2024, 2, 1, 12), Date(2024, 2, 1, 11), 2, 1, 12))
	this.False(Yearly(Date(2024, 2, 1, 12), Date(2024, 2, 1, 12), 2, 1, 12))
	this.False(Yearly(Date(2024, 2, 1, 12), Date(2024, 3, 1, 12), 2, 1, 12))
	this.False(Yearly(Date(2024, 2, 1, 12), Date(2023, 1, 1, 12), 2, 31, 12))
}

func (this *SuiteTime) TestYearlyPrev() {
	this.Equal(Date(2022, 2, 1, 12), YearlyPrev(Date(2023, 2, 1, 11), 2, 1, 12))
	this.Equal(Date(2023, 2, 1, 12), YearlyPrev(Date(2023, 2, 1, 12), 2, 1, 12))
	this.Equal(Date(2023, 2, 1, 12), YearlyPrev(Date(2023, 2, 1, 13), 2, 1, 12))
	this.Equal(Date(2023, 3, 31, 12), YearlyPrev(Date(2023, 4, 1, 13), 3, 31, 12))
	this.Equal(Date(2023, 1, 1, 12), YearlyPrev(Date(2023, 2, 1, 13), -1, 1, 12))
	this.Equal(Date(2022, 12, 1, 12), YearlyPrev(Date(2023, 3, 1, 13), 99, 1, 12))
	this.Equal(Date(2023, 1, 1, 12), YearlyPrev(Date(2023, 2, 1, 13), 1, -1, 12))
	this.Equal(Date(2023, 1, 31, 12), YearlyPrev(Date(2023, 2, 1, 13), 1, 99, 12))
	this.Equal(Date(2023, 1, 1, 0), YearlyPrev(Date(2023, 2, 1, 1), 1, 1, -1))
	this.Equal(Date(2023, 1, 1, 23), YearlyPrev(Date(2023, 2, 1, 0), 1, 1, 99))
	this.Equal(time.Time{}, YearlyPrev(Date(2023, 2, 1, 0), 2, 31, 12))
}

func (this *SuiteTime) TestYearlyNext() {
	this.Equal(Date(2023, 2, 1, 12), YearlyNext(Date(2023, 2, 1, 11), 2, 1, 12))
	this.Equal(Date(2024, 2, 1, 12), YearlyNext(Date(2023, 2, 1, 12), 2, 1, 12))
	this.Equal(Date(2024, 2, 1, 12), YearlyNext(Date(2023, 2, 1, 13), 2, 1, 12))
	this.Equal(Date(2024, 1, 1, 12), YearlyNext(Date(2023, 2, 1, 13), -1, 1, 12))
	this.Equal(Date(2023, 12, 1, 12), YearlyNext(Date(2023, 3, 1, 13), 99, 1, 12))
	this.Equal(Date(2024, 1, 1, 12), YearlyNext(Date(2023, 2, 1, 13), 1, -1, 12))
	this.Equal(Date(2024, 1, 31, 12), YearlyNext(Date(2023, 2, 1, 13), 1, 99, 12))
	this.Equal(Date(2024, 1, 1, 0), YearlyNext(Date(2023, 2, 1, 1), 1, 1, -1))
	this.Equal(Date(2024, 1, 1, 23), YearlyNext(Date(2023, 2, 1, 0), 1, 1, 99))
	this.Equal(time.Time{}, YearlyNext(Date(2023, 2, 1, 0), 2, 31, 12))
}

func (this *SuiteTime) TestFixedPrev() {
	this.Equal(Date(1970, 1, 1), FixedPrev(time.Time{}, time.Time{}, TimeHour))
	this.Equal(Date(2023, 1, 1, 8), FixedPrev(time.Time{}, Date(2023, 1, 1, 8), TimeHour))
	this.Equal(Date(2023, 1, 1, 8), FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8), TimeHour))
	this.Equal(Date(2023, 1, 1, 8), FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 10), TimeHour))
	this.Equal(Date(2023, 1, 1, 8), FixedPrev(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 50), TimeHour))
	this.Equal(Date(2023, 1, 1, 7), FixedPrev(Date(1970, 1, 1, 1), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(Date(2023, 1, 1, 8), FixedPrev(Date(1970, 1, 1, 2), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(Date(2023, 1, 1, 6, 30), FixedPrev(Date(1970, 1, 1, 3, 30), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(time.Time{}, FixedPrev(time.Time{}, time.Time{}, 0))
}

func (this *SuiteTime) TestFixedNext() {
	this.Equal(Date(1970, 1, 1, 1), FixedNext(time.Time{}, time.Time{}, TimeHour))
	this.Equal(Date(2023, 1, 1, 9), FixedNext(time.Time{}, Date(2023, 1, 1, 8), TimeHour))
	this.Equal(Date(2023, 1, 1, 9), FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8), TimeHour))
	this.Equal(Date(2023, 1, 1, 9), FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 10), TimeHour))
	this.Equal(Date(2023, 1, 1, 9), FixedNext(Date(1970, 1, 1, 8), Date(2023, 1, 1, 8, 50), TimeHour))
	this.Equal(Date(2023, 1, 1, 10), FixedNext(Date(1970, 1, 1, 1), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(Date(2023, 1, 1, 11), FixedNext(Date(1970, 1, 1, 2), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(Date(2023, 1, 1, 9, 30), FixedNext(Date(1970, 1, 1, 3, 30), Date(2023, 1, 1, 8, 1, 0, 0), TimeHour*3))
	this.Equal(time.Time{}, FixedNext(time.Time{}, time.Time{}, 0))
}

func (this *SuiteTime) TestCalculateDays() {
	this.Equal(1, CalculateDays(Date(2023, 2, 10), Date(2023, 2, 11)))
	this.Equal(1, CalculateDays(Date(2023, 2, 10, 1), Date(2023, 2, 11, 1)))
	this.Equal(2, CalculateDays(Date(2023, 2, 10, 1), Date(2023, 2, 12, 1)))
	this.Equal(3, CalculateDays(Date(2023, 2, 10, 1), Date(2023, 2, 13, 1)))
	this.Equal(1, CalculateDays(Date(2023, 2, 11), Date(2023, 2, 10)))
	this.Equal(1, CalculateDays(Date(2023, 2, 11, 1), Date(2023, 2, 10, 1)))
	this.Equal(2, CalculateDays(Date(2023, 2, 12, 1), Date(2023, 2, 10, 1)))
	this.Equal(3, CalculateDays(Date(2023, 2, 13, 1), Date(2023, 2, 10, 1)))
}

func (this *SuiteTime) TestCalculateDaysWithBaseline() {
	this.Equal(1, CalculateDaysWithBaseline(Date(2023, 2, 10), Date(2023, 2, 11), TimeHour*5))
	this.Equal(1, CalculateDaysWithBaseline(Date(2023, 2, 10, 1), Date(2023, 2, 11, 1), TimeHour*5))
	this.Equal(2, CalculateDaysWithBaseline(Date(2023, 2, 10, 1), Date(2023, 2, 12, 1), TimeHour*5))
	this.Equal(3, CalculateDaysWithBaseline(Date(2023, 2, 10, 1), Date(2023, 2, 13, 1), TimeHour*5))
}

func (this *SuiteTime) TestClampHour() {
	this.Equal(0, clampHour(-1))
	this.Equal(23, clampHour(99))
}

func (this *SuiteTime) TestClampWeek() {
	this.Equal(6, clampWeek(-1))
	this.Equal(0, clampWeek(7))
}

func (this *SuiteTime) TestClampMDay() {
	this.Equal(1, clampMDay(-1))
	this.Equal(31, clampMDay(99))
}

func (this *SuiteTime) TestClampMonth() {
	this.Equal(1, clampMonth(-1))
	this.Equal(12, clampMonth(99))
}
