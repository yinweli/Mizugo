package helps

import (
	"fmt"
	"time"
)

const (
	TimeMillisecond = time.Millisecond      // 1毫秒時間
	TimeSecond      = time.Second           // 1秒時間
	TimeMinute      = time.Minute           // 1分鐘時間
	TimeHour        = time.Hour             // 1小時時間
	TimeDay         = TimeHour * 24         // 1日時間
	TimeWeek        = TimeDay * 7           // 1週時間
	TimeMonth       = TimeWeek * 4          // 1月時間, 這個時間並不準確
	TimeYear        = TimeDay * 365         // 1年時間, 這個時間並不準確
	LayoutSecond    = "2006-01-02 15:04:05" // 時間字串格式, 使用年月日時分秒
	LayoutMinute    = "2006-01-02 15:04"    // 時間字串格式, 使用年月日時分
	LayoutDay       = "2006-01-02"          // 時間字串格式, 使用年月日
	firstYear       = 1970                  // unix計時的第一年
)

// SetTimeZone 設定時區資料
func SetTimeZone(timeZone string) error {
	location, err := time.LoadLocation(timeZone)

	if err != nil {
		return fmt.Errorf("setTimeZone: %w", err)
	} // if

	zone = location
	return nil
}

// SetTimeZoneUTC 設定時區資料為UTC時間
func SetTimeZoneUTC() {
	zone = time.UTC
}

// SetTimeZoneLocal 設定時區資料為本地時間
func SetTimeZoneLocal() {
	zone = time.Local
}

// GetTimeZone 取得時區資料
func GetTimeZone() *time.Location {
	if zone == nil {
		zone = time.UTC
	} // if

	return zone
}

// Time 取得現在時間, 會轉換為 SetTimeZone 設定的時區時間
func Time() time.Time {
	return time.Now().In(GetTimeZone())
}

// Timef 取得格式時間, 會轉換為 SetTimeZone 設定的時區時間, 時間字串按照 layout 來解析
func Timef(layout, v string) (time.Time, error) {
	if v == "" {
		return time.Time{}, nil
	} // if

	u, err := time.ParseInLocation(layout, v, GetTimeZone())

	if err != nil {
		return time.Time{}, fmt.Errorf("timef: %w", err)
	} // if

	return u, nil
}

// Date 取得指定時間, 會轉換為 SetTimeZone 設定的時區時間;
// 輸入參數依序是 年, 月, 日, 時, 分, 秒, 毫秒; 若未輸入則自動填0;
// 例如 Date(2023, 2, 15) 會得到 2023/02/15 00:00:00 的時間
func Date(v ...int) time.Time {
	if s := len(v); s >= 7 { //nolint:gomnd
		return time.Date(v[0], time.Month(v[1]), v[2], v[3], v[4], v[5], v[6], GetTimeZone())
	} else if s >= 6 { //nolint:gomnd
		return time.Date(v[0], time.Month(v[1]), v[2], v[3], v[4], v[5], 0, GetTimeZone())
	} else if s >= 5 { //nolint:gomnd
		return time.Date(v[0], time.Month(v[1]), v[2], v[3], v[4], 0, 0, GetTimeZone())
	} else if s >= 4 { //nolint:gomnd
		return time.Date(v[0], time.Month(v[1]), v[2], v[3], 0, 0, 0, GetTimeZone())
	} else if s >= 3 {
		return time.Date(v[0], time.Month(v[1]), v[2], 0, 0, 0, 0, GetTimeZone())
	} else if s >= 2 {
		return time.Date(v[0], time.Month(v[1]), 1, 0, 0, 0, 0, GetTimeZone())
	} else if s >= 1 {
		return time.Date(v[0], 1, 1, 0, 0, 0, 0, GetTimeZone())
	} else {
		return time.Time{}
	} // if
}

// Between 檢查 u 是否在 s 與 e 之間, 當 s 與 e 為空時, 回傳true
func Between(s, e, u time.Time) bool {
	if s.IsZero() && e.IsZero() {
		return true
	} // if

	if s.IsZero() {
		return u.Before(e)
	} // if

	if e.IsZero() {
		return u.After(s)
	} // if

	return u.After(s) && u.Before(e)
}

// Overlap 檢查兩個時間段是否有重疊
func Overlap(s1, e1, s2, e2 time.Time) bool {
	return e1.After(s2) && e2.After(s1)
}

// Daily 檢查每日是否到期
func Daily(now, last time.Time, hour int) bool {
	next := DailyNext(last, hour)
	return now.Equal(next) || now.After(next)
}

// DailyPrev 取得上次的每日時間
func DailyPrev(now time.Time, hour int) time.Time {
	prev := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Before(prev) { // 如果已經經過時間, 就設定為昨天
		prev = prev.AddDate(0, 0, -1)
	} // if

	return prev
}

// DailyNext 取得下次的每日時間
func DailyNext(now time.Time, hour int) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為明天
		next = next.AddDate(0, 0, 1)
	} // if

	return next
}

// Weekly 檢查每週是否到期, wday是 time.Weekday, 週日為0, 週一為1, 依此類推
func Weekly(now, last time.Time, wday, hour int) bool {
	next := WeeklyNext(last, wday, hour)
	return now.Equal(next) || now.After(next)
}

// WeeklyPrev 取得上次的每週時間, wday是 time.Weekday, 週日為0, 週一為1, 依此類推
func WeeklyPrev(now time.Time, wday, hour int) time.Time {
	day := int(time.Weekday(wday) - now.Weekday())

	if day > 0 { // 如果還沒到時間, 就設定為上週
		day -= 7
	} // if

	if day == 0 && now.Hour() < hour { // 如果還沒到時間, 就設定為上週
		day -= 7
	} // if

	prev := now.AddDate(0, 0, day)
	return time.Date(prev.Year(), prev.Month(), prev.Day(), hour, 0, 0, 0, now.Location())
}

// WeeklyNext 取得下次的每週時間, wday是 time.Weekday, 週日為0, 週一為1, 依此類推
func WeeklyNext(now time.Time, wday, hour int) time.Time {
	day := int(time.Weekday(wday) - now.Weekday())

	if day < 0 { // 如果已經經過時間, 就設定為下週
		day += 7
	} // if

	if day == 0 && now.Hour() >= hour { // 如果已經經過時間, 就設定為下週
		day += 7
	} // if

	next := now.AddDate(0, 0, day)
	return time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, now.Location())
}

// Monthly 檢查每月是否到期, mday是month-day
func Monthly(now, last time.Time, mday, hour int) bool {
	next := MonthlyNext(last, mday, hour)
	return now.Equal(next) || now.After(next)
}

// MonthlyPrev 取得上次的每月時間, mday是month-day
func MonthlyPrev(now time.Time, mday, hour int) time.Time {
	prev := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Before(prev) || now.Hour() < hour { // 如果還沒到時間，就設定為上個月
		prev = prev.AddDate(0, -1, 0)
	} // if

	// 當月份日數不同時, 例如輸入的日數為31, 但是當月日數最多只到28時, 就會將日期減1, 直到日期有效為止
	for prev.Month() == now.Month() && prev.Day() != mday {
		prev = prev.AddDate(0, 0, -1)
	} // if

	return prev
}

// MonthlyNext 取得下次的每月時間, mday是month-day
func MonthlyNext(now time.Time, mday, hour int) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為下個月
		next = next.AddDate(0, 1, 0)
	} // if

	return next
}

// FixedPrev 取得上個固定間隔時間;
// base 若小於1970年, base 會被改為從1970年開始;
// now 若小於 base, now 會被改為 base;
// duration 不能小於1秒, 會造成計算錯誤
func FixedPrev(base, now time.Time, duration time.Duration) time.Time {
	if base.Before(time.Date(firstYear, 1, 1, 0, 0, 0, 0, base.Location())) {
		base = time.Date(firstYear, base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), 0, base.Location())
	} // if

	if now.Before(base) {
		now = base
	} // if

	nsec := int64(now.Sub(base).Seconds())
	dsec := int64(duration.Seconds())
	quotient := nsec / dsec

	if quotient < 0 {
		quotient = 0
	} // if

	return base.Add(time.Duration(quotient*dsec) * time.Second)
}

// FixedNext 取得下個固定間隔時間;
// base 若小於1970年, base 會被改為從1970年開始;
// now 若小於 base, now 會被改為 base;
// duration 不能小於1秒, 會造成計算錯誤
func FixedNext(base, now time.Time, duration time.Duration) time.Time {
	if base.Before(time.Date(firstYear, 1, 1, 0, 0, 0, 0, base.Location())) {
		base = time.Date(firstYear, base.Month(), base.Day(), base.Hour(), base.Minute(), base.Second(), 0, base.Location())
	} // if

	if now.Before(base) {
		now = base
	} // if

	nsec := int64(now.Sub(base).Seconds())
	dsec := int64(duration.Seconds())
	quotient := nsec / dsec

	if quotient < 0 {
		quotient = 1
	} else {
		quotient++
	} // if

	return base.Add(time.Duration(quotient*dsec) * time.Second)
}

var zone *time.Location // 時區資料
