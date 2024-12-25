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
	TimeDay         = TimeHour * DayHourMax // 1日時間
	TimeWeek        = TimeDay * WeekdayMax  // 1週時間
	TimeMonth       = TimeWeek * 4          // 1月時間, 這個時間並不準確
	TimeYear        = TimeDay * 365         // 1年時間, 這個時間並不準確
	DayHourMax      = 24                    // 一天有幾小時
	WeekdayMax      = 7                     // 1週有幾天
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

// Timef 取得格式時間, 會轉換為 SetTimeZone 設定的時區時間, value 按照 layout 來解析
func Timef(layout, value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	} // if

	u, err := time.ParseInLocation(layout, value, GetTimeZone())

	if err != nil {
		return time.Time{}, fmt.Errorf("timef: %w", err)
	} // if

	return u, nil
}

// Date 取得指定時間, 會轉換為 SetTimeZone 設定的時區時間;
// 輸入參數依序是 年, 月, 日, 時, 分, 秒, 毫秒; 若時, 分, 秒, 毫秒未輸入則自動填0;
// 例如 Date(2023, 2, 15) 會得到 2023/02/15 00:00:00 的時間
func Date(year int, month time.Month, day int, value ...int) time.Time {
	hour, minute, sec, nsec := 0, 0, 0, 0

	if len(value) > 0 {
		hour = value[0]
	} // if

	if len(value) > 1 {
		minute = value[1]
	} // if

	if len(value) > 2 {
		sec = value[2]
	} // if

	if len(value) > 3 {
		nsec = value[3]
	} // if

	return time.Date(year, month, day, hour, minute, sec, nsec, GetTimeZone())
}

// Before 檢查 u 是否在 t 之前
func Before(u, t time.Time) bool {
	return u.Before(t)
}

// Beforef 檢查 u 是否在 t 之前, u 按照 layout 來解析
func Beforef(layout, u string, t time.Time) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	return uf.Before(t)
}

// Beforefx 檢查 u 是否在 t 之前, u 與 t 按照 layout 來解析
func Beforefx(layout, u, t string) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	tf, err := Timef(layout, t)

	if err != nil {
		return false
	} // if

	return uf.Before(tf)
}

// After 檢查 u 是否在 t 之後
func After(u, t time.Time) bool {
	return u.After(t)
}

// Afterf 檢查 u 是否在 t 之後, u 按照 layout 來解析
func Afterf(layout, u string, t time.Time) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	return uf.After(t)
}

// Afterfx 檢查 u 是否在 t 之後, u 與 t 按照 layout 來解析
func Afterfx(layout, u, t string) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	tf, err := Timef(layout, t)

	if err != nil {
		return false
	} // if

	return uf.After(tf)
}

// Between 檢查 u 是否在 start 與 end 之間;
// 當 start 與 end 為空時, 回傳 zero, zero 預設為true
func Between(start, end, u time.Time, zero ...bool) bool {
	if start.IsZero() && end.IsZero() {
		if len(zero) > 0 {
			return zero[0]
		} else {
			return true
		} // if
	} // if

	if start.IsZero() {
		return u.Before(end)
	} // if

	if end.IsZero() {
		return u.After(start)
	} // if

	return u.After(start) && u.Before(end)
}

// Betweenf 檢查 u 是否在 start 與 end 之間, start 與 end 按照 layout 來解析;
// 當 start 與 end 為空時, 回傳 zero, zero 預設為true
func Betweenf(layout, start, end string, u time.Time, zero ...bool) bool {
	startf, err := Timef(layout, start)

	if err != nil {
		return false
	} // if

	endf, err := Timef(layout, end)

	if err != nil {
		return false
	} // if

	return Between(startf, endf, u, zero...)
}

// Betweenfx 檢查 u 是否在 start 與 end 之間, start 與 end 與 u 按照 layout 來解析;
// 當 start 與 end 為空時, 回傳 zero, zero 預設為true
func Betweenfx(layout, start, end, u string, zero ...bool) bool {
	startf, err := Timef(layout, start)

	if err != nil {
		return false
	} // if

	endf, err := Timef(layout, end)

	if err != nil {
		return false
	} // if

	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	return Between(startf, endf, uf, zero...)
}

// Overlap 檢查兩個時間段是否有重疊
func Overlap(start1, end1, start2, end2 time.Time) bool {
	return end1.After(start2) && end2.After(start1)
}

// Overlapf 檢查兩個時間段是否有重疊, start1 與 end1 按照 layout 來解析
func Overlapf(layout, start1, end1 string, start2, end2 time.Time) bool {
	start1f, err := Timef(layout, start1)

	if err != nil {
		return false
	} // if

	end1f, err := Timef(layout, end1)

	if err != nil {
		return false
	} // if

	return Overlap(start1f, end1f, start2, end2)
}

// Overlapfx 檢查兩個時間段是否有重疊, start1 與 end1 與 start2 與 end2 按照 layout 來解析
func Overlapfx(layout, start1, end1, start2, end2 string) bool {
	start1f, err := Timef(layout, start1)

	if err != nil {
		return false
	} // if

	end1f, err := Timef(layout, end1)

	if err != nil {
		return false
	} // if

	start2f, err := Timef(layout, start2)

	if err != nil {
		return false
	} // if

	end2f, err := Timef(layout, end2)

	if err != nil {
		return false
	} // if

	return Overlap(start1f, end1f, start2f, end2f)
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
	prev := time.Date(now.Year(), now.Month(), mday, hour, 0, 0, 0, now.Location())

	if now.Before(prev) { // 如果還沒到時間，就設定為上個月
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
	next := time.Date(now.Year(), now.Month(), mday, hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為下個月
		next = next.AddDate(0, 1, 0)
	} // if

	return next
}

// Yearly 檢查每年是否到期, mday是month-day
func Yearly(now, last time.Time, month, mday, hour int) bool {
	next := YearlyNext(last, month, mday, hour)
	return now.Equal(next) || now.After(next)
}

// YearlyPrev 取得上次的每年時間, mday是month-day
func YearlyPrev(now time.Time, month, mday, hour int) time.Time {
	prev := time.Date(now.Year(), time.Month(month), mday, hour, 0, 0, 0, now.Location())

	// 當月份日數不同時, 例如輸入的日數為31, 但是當月日數最多只到28時, 就會將日期減1, 直到日期有效為止
	for prev.Month() == now.Month() && prev.Day() != mday {
		prev = prev.AddDate(0, 0, -1)
	} // if

	if now.Before(prev) { // 如果還沒到時間，就設定為去年
		prev = prev.AddDate(-1, 0, 0)
	} // if

	return prev
}

// YearlyNext 取得下次的每年時間, mday是month-day
func YearlyNext(now time.Time, month, mday, hour int) time.Time {
	next := time.Date(now.Year(), time.Month(month), mday, hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為明年
		next = next.AddDate(1, 0, 0)
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

// CalculateDays 計算兩個時間相差幾天
func CalculateDays(t1, t2 time.Time) int {
	if t2.Before(t1) {
		return CalculateDays(t2, t1)
	} // if

	t1x := Date(t1.Year(), t1.Month(), t1.Day())
	t2x := Date(t2.Year(), t2.Month(), t2.Day())
	return int(t2x.Sub(t1x).Hours() / DayHourMax)
}

// CalculateDaysWithBaseline 計算兩個時間相差幾天, 但是以 base 指定的時間為日界
// base 應該是個介於 0 ~ 23:59:59 的時間長度, 例如 base 為5小時的話, 表示以凌晨5點為日界
func CalculateDaysWithBaseline(t1, t2 time.Time, base time.Duration) int {
	return CalculateDays(t1.Add(-base), t2.Add(-base))
}

var zone *time.Location // 時區資料
