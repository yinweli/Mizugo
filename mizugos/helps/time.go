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
	TimeMonth       = TimeWeek * 4          // 1月時間
	TimeYear        = TimeDay * 365         // 1年時間
	LayoutSecond    = "2006-01-02 15:04:05" // 時間字串格式, 使用年月日時分秒
	LayoutMinute    = "2006-01-02 15:04"    // 時間字串格式, 使用年月日時分
)

// SetTimeZone 設定時區資料
func SetTimeZone(timeZone string) error {
	zone, err := time.LoadLocation(timeZone)

	if err != nil {
		return fmt.Errorf("setTimeZone: %w", err)
	} // if

	location = zone
	return nil
}

// GetTimeZone 取得時區資料
func GetTimeZone() *time.Location {
	if location == nil {
		location = time.UTC
	} // if

	return location
}

// Time 取得現在時間, 會轉換為 SetTimeZone 設定的時區時間, 預設是UTC+0
func Time() time.Time {
	return time.Now().In(GetTimeZone())
}

// Date 取得指定時間, 會轉換為 SetTimeZone 設定的時區時間, 預設是UTC+0
func Date(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, GetTimeZone())
}

// Between 檢查時間是否在開始與結束時間之間, 當開始與結束時間為空時, 會回傳預設值
func Between(start, end, now time.Time, preset bool) bool {
	switch {
	case start.IsZero() && end.IsZero():
		return preset

	case start.IsZero() && end.IsZero() == false:
		return end.After(now)

	case start.IsZero() == false && end.IsZero():
		return start.Before(now)

	default:
		return start.Before(now) && end.After(now)
	} // switch
}

// Betweenf 檢查時間是否在開始與結束時間字串之間, 當開始與結束時間字串為空時, 會回傳預設值, 時間字串按照 layout 來解析
func Betweenf(layout, start, end string, now time.Time, preset bool) bool {
	startTime := time.Time{}
	endTime := time.Time{}
	err := error(nil)

	if start != "" {
		if startTime, err = time.ParseInLocation(layout, start, GetTimeZone()); err != nil {
			return false
		} // if
	} // if

	if end != "" {
		if endTime, err = time.ParseInLocation(layout, end, GetTimeZone()); err != nil {
			return false
		} // if
	} // if

	return Between(startTime, endTime, now, preset)
}

// SameDay 檢查兩個時間是否為同一天
func SameDay(t1, t2 time.Time) bool {
	t1 = t1.In(location) // 轉換到相同時區, 避免時區問題
	t2 = t2.In(location)
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

// Daily 檢查每日是否到期
func Daily(now, last time.Time, hour int) bool {
	next := DailyNext(last, hour)
	return now.Equal(next) || now.After(next)
}

// DailyPrev 取得上次的每日時間
func DailyPrev(now time.Time, hour int) time.Time {
	prev := Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0)

	if now.Before(prev) { // 如果已經經過時間, 就設定為昨天
		prev = prev.AddDate(0, 0, -1)
	} // if

	return prev
}

// DailyNext 取得下次的每日時間
func DailyNext(now time.Time, hour int) time.Time {
	next := Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0)

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
	return Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0)
}

// Monthly 檢查每月是否到期, mday是month-day
func Monthly(now, last time.Time, mday, hour int) bool {
	next := MonthlyNext(last, mday, hour)
	return now.Equal(next) || now.After(next)
}

// MonthlyPrev 取得上次的每月時間, mday是month-day
func MonthlyPrev(now time.Time, mday, hour int) time.Time {
	prev := Date(now.Year(), now.Month(), mday, hour, 0, 0, 0)

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
	next := Date(now.Year(), now.Month(), mday, hour, 0, 0, 0)

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為下個月
		next = next.AddDate(0, 1, 0)
	} // if

	return next
}

// FixedPrev 取得上個固定間隔時間; duration 不能小於1秒, 並且 86400秒 % duration必須為0, 否則會造成計算誤差, 可以用 FixedCheck 檢查 duration 是否正確
func FixedPrev(now time.Time, duration time.Duration) time.Time {
	nsec := now.Unix()
	dsec := int64(duration.Seconds())
	quotient := nsec / dsec
	return time.Unix(quotient*dsec, 0).In(GetTimeZone())
}

// FixedNext 取得下個固定間隔時間; duration 不能小於1秒, 並且 86400秒 % duration必須為0, 否則會造成計算誤差, 可以用 FixedCheck 檢查 duration 是否正確
func FixedNext(now time.Time, duration time.Duration) time.Time {
	nsec := now.Unix()
	dsec := int64(duration.Seconds())
	quotient := nsec/dsec + 1
	return time.Unix(quotient*dsec, 0).In(GetTimeZone())
}

// FixedCheck 檢查提供給 FixedPrev 與 FixedNext 的 duration 是否正確
func FixedCheck(duration time.Duration) bool {
	return TimeDay%duration == 0
}

var location *time.Location // 時區資料
