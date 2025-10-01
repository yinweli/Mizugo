package helps

import (
	"fmt"
	"time"
)

const (
	TimeMillisecond = time.Millisecond      // 1 毫秒時間單位
	TimeSecond      = time.Second           // 1 秒時間單位
	TimeMinute      = time.Minute           // 1 分鐘時間單位
	TimeHour        = time.Hour             // 1 小時時間單位
	TimeDay         = TimeHour * DayHourMax // 1 日時間單位
	TimeWeek        = TimeDay * WeekdayMax  // 1 週時間單位
	TimeMonth       = TimeWeek * 4          // 1 月時間單位(4 週)
	TimeYear        = TimeDay * 365         // 1 年時間單位(365 天)
	DayHourMax      = 24                    // 1 天有幾小時
	WeekdayMax      = 7                     // 1 週有幾天
	LayoutSecond    = "2006-01-02 15:04:05" // 時間格式 YYYY-MM-DD HH:mm:SS
	LayoutMinute    = "2006-01-02 15:04"    // 時間格式 YYYY-MM-DD HH:mm
	LayoutDay       = "2006-01-02"          // 時間格式 YYYY-MM-DD
)

// SetTimeZone 設定時區
//
// 範例:
//
//	SetTimeZone("UTC")
//	SetTimeZone("Asia/Taipei")
//	SetTimeZone("America/Los_Angeles")
func SetTimeZone(timeZone string) error {
	location, err := time.LoadLocation(timeZone)

	if err != nil {
		return fmt.Errorf("setTimeZone: %w", err)
	} // if

	zone = location
	return nil
}

// SetTimeZoneUTC 將時區設為 UTC
func SetTimeZoneUTC() {
	zone = time.UTC
}

// SetTimeZoneLocal 將時區設為本機時區
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

// Time 取得現在時間, 並套用 SetTimeZone 設定的時區時間
func Time() time.Time {
	return time.Now().In(GetTimeZone())
}

// Timef 取得依 layout 解析時間, 並套用 SetTimeZone 設定的時區時間
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

// Date 取得指定時間, 並套用 SetTimeZone 設定的時區時間
//
// 輸入參數依序是 年, 月, 日, 時 (可省略), 分 (可省略), 秒 (可省略), 毫秒 (可省略)
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

// Beforef 檢查依 layout 解析後的時間 u 是否在 t 之前
func Beforef(layout, u string, t time.Time) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	return uf.Before(t)
}

// Beforefx 檢查依 layout 解析後的時間 u 是否在 t 之前
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

// Afterf 檢查依 layout 解析後的時間 u 是否在 t 之後
func Afterf(layout, u string, t time.Time) bool {
	uf, err := Timef(layout, u)

	if err != nil {
		return false
	} // if

	return uf.After(t)
}

// Afterfx 檢查依 layout 解析後的時間 u 是否在 t 之後
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

// Between 檢查 u 是否在 start 與 end 之間, 當 start 與 end 為空時, 回傳 zero 或是預設 true
func Between(start, end, u time.Time, zero ...bool) bool {
	if start.IsZero() && end.IsZero() {
		if len(zero) > 0 {
			return zero[0]
		} // if

		return true
	} // if

	if start.IsZero() {
		return u.Before(end)
	} // if

	if end.IsZero() {
		return u.After(start)
	} // if

	return u.After(start) && u.Before(end)
}

// Betweenf 檢查依 layout 解析後的時間 u 是否在 start 與 end 之間, 當 start 與 end 為空時, 回傳 zero 或是預設 true
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

// Betweenfx 檢查依 layout 解析後的時間 u 是否在 start 與 end 之間, 當 start 與 end 為空時, 回傳 zero 或是預設 true
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

// Overlapf 檢查依 layout 解析後的兩個時間段是否有重疊
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

// Overlapfx 檢查依 layout 解析後的兩個時間段是否有重疊
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

// Daily 檢查每日時間是否已到期/逾期
func Daily(now, last time.Time, hour int) bool {
	next := DailyNext(last, hour)
	return now.Equal(next) || now.After(next)
}

// DailyPrev 取得上一個每日時間
func DailyPrev(now time.Time, hour int) time.Time {
	hour = clampHour(hour)
	prev := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Before(prev) { // 如果尚未到達時間, 就設定為昨天
		prev = prev.AddDate(0, 0, -1)
	} // if

	return prev
}

// DailyNext 取得下一個每日時間
func DailyNext(now time.Time, hour int) time.Time {
	hour = clampHour(hour)
	next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為明天
		next = next.AddDate(0, 0, 1)
	} // if

	return next
}

// Weekly 檢查每週時間是否已到期/逾期
//
// wday = 0 為週日, wday = 1 為週一, wday = 6 為週六
func Weekly(now, last time.Time, wday, hour int) bool {
	next := WeeklyNext(last, wday, hour)
	return now.Equal(next) || now.After(next)
}

// WeeklyPrev 取得上一個每週時間
//
// wday = 0 為週日, wday = 1 為週一, wday = 6 為週六
func WeeklyPrev(now time.Time, wday, hour int) time.Time {
	wday = clampWeek(wday)
	hour = clampHour(hour)
	day := wday - int(now.Weekday())

	if day > 0 { // 如果尚未到達時間, 就設定為上週
		day -= 7
	} // if

	if day == 0 && now.Hour() < hour { // 如果尚未到達時間, 就設定為上週
		day -= 7
	} // if

	prev := now.AddDate(0, 0, day)
	return time.Date(prev.Year(), prev.Month(), prev.Day(), hour, 0, 0, 0, now.Location())
}

// WeeklyNext 取得下一個每週時間
//
// wday = 0 為週日, wday = 1 為週一, wday = 6 為週六
func WeeklyNext(now time.Time, wday, hour int) time.Time {
	wday = clampWeek(wday)
	hour = clampHour(hour)
	day := wday - int(now.Weekday())

	if day < 0 { // 如果已經經過時間, 就設定為下週
		day += 7
	} // if

	if day == 0 && now.Hour() >= hour { // 如果已經經過時間, 就設定為下週
		day += 7
	} // if

	next := now.AddDate(0, 0, day)
	return time.Date(next.Year(), next.Month(), next.Day(), hour, 0, 0, 0, now.Location())
}

// Monthly 檢查每月時間是否已到期/逾期
func Monthly(now, last time.Time, mday, hour int) bool {
	next := MonthlyNext(last, mday, hour)
	return now.Equal(next) || now.After(next)
}

// MonthlyPrev 取得上一個每月時間
//
// 若指定日期錯誤(例: 輸入 31 日, 但該月僅有 28 日), 則會逐日遞減, 直到落在有效日期為止
func MonthlyPrev(now time.Time, mday, hour int) time.Time {
	mday = clampMDay(mday)
	hour = clampHour(hour)
	prev := time.Date(now.Year(), now.Month(), mday, hour, 0, 0, 0, now.Location())

	if now.Before(prev) { // 如果尚未到達時間，就設定為上個月
		prev = prev.AddDate(0, -1, 0)
	} // if

	for prev.Day() != mday { // 若指定日期錯誤, 則會逐日遞減, 直到落在有效日期為止
		prev = prev.AddDate(0, 0, -1)
	} // if

	return prev
}

// MonthlyNext 取得下一個每月時間
//
// 若指定日期錯誤(例: 輸入 31 日, 但該月僅有 28 日), 則會逐日遞增, 直到落在有效日期為止
func MonthlyNext(now time.Time, mday, hour int) time.Time {
	mday = clampMDay(mday)
	hour = clampHour(hour)
	next := time.Date(now.Year(), now.Month(), mday, hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為下個月
		next = next.AddDate(0, 1, 0)
	} // if

	for next.Day() != mday { // 若指定日期錯誤, 則會逐日遞增, 直到落在有效日期為止
		next = next.AddDate(0, 0, 1)
	} // if

	return next
}

// Yearly 檢查每年時間是否已到期/逾期
func Yearly(now, last time.Time, month, mday, hour int) bool {
	next := YearlyNext(last, month, mday, hour)
	return next.IsZero() == false && (now.Equal(next) || now.After(next))
}

// YearlyPrev 取得上一個每年時間
//
// 若指定日期錯誤(例: 輸入 31 日, 但該月僅有 28 日), 回傳 time.Time{} 表示錯誤
func YearlyPrev(now time.Time, month, mday, hour int) time.Time {
	month = clampMonth(month)
	mday = clampMDay(mday)
	hour = clampHour(hour)
	prev := time.Date(now.Year(), time.Month(month), mday, hour, 0, 0, 0, now.Location())

	if now.Before(prev) { // 如果尚未到達時間，就設定為去年
		prev = prev.AddDate(-1, 0, 0)
	} // if

	if prev.Month() != time.Month(month) || prev.Day() != mday { // 若指定日期錯誤, 回傳 time.Time{} 表示錯誤
		return time.Time{}
	} // if

	return prev
}

// YearlyNext 取得下一個每年時間
//
// 若指定日期錯誤(例: 輸入 31 日, 但該月僅有 28 日), 回傳 time.Time{} 表示錯誤
func YearlyNext(now time.Time, month, mday, hour int) time.Time {
	month = clampMonth(month)
	mday = clampMDay(mday)
	hour = clampHour(hour)
	next := time.Date(now.Year(), time.Month(month), mday, hour, 0, 0, 0, now.Location())

	if now.Equal(next) || now.After(next) { // 如果已經經過時間, 就設定為明年
		next = next.AddDate(1, 0, 0)
	} // if

	if next.Month() != time.Month(month) || next.Day() != mday { // 若指定日期錯誤, 回傳 time.Time{} 表示錯誤
		return time.Time{}
	} // if

	return next
}

// FixedPrev 取得上個固定間隔時間
//   - 若 base < 1970-01-01, 則以 1970-01-01 00:00:00 作為新的 base
//   - 若 now < base, 則 now = base
//   - duration 必須 > 0, 否則回傳 time.Time{}
func FixedPrev(base, now time.Time, duration time.Duration) time.Time {
	if duration <= 0 {
		return time.Time{}
	} // if

	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, base.Location())

	if base.Before(epoch) {
		base = epoch
	} // if

	if now.Before(base) {
		now = base
	} // if

	nsec := time.Duration(now.Sub(base).Seconds())
	dsec := time.Duration(duration.Seconds())
	quotient := nsec / dsec

	if quotient < 0 {
		quotient = 0
	} // if

	return base.Add(quotient * dsec * time.Second)
}

// FixedNext 取得下個固定間隔時間
//   - 若 base < 1970-01-01, 則以 1970-01-01 00:00:00 作為新的 base
//   - 若 now < base, 則 now = base
//   - duration 必須 > 0, 否則回傳 time.Time{}
func FixedNext(base, now time.Time, duration time.Duration) time.Time {
	if duration <= 0 {
		return time.Time{}
	} // if

	epoch := time.Date(1970, 1, 1, 0, 0, 0, 0, base.Location())

	if base.Before(epoch) {
		base = epoch
	} // if

	if now.Before(base) {
		now = base
	} // if

	nsec := time.Duration(now.Sub(base).Seconds())
	dsec := time.Duration(duration.Seconds())
	quotient := nsec / dsec

	if quotient < 0 {
		quotient = 1
	} else {
		quotient++
	} // if

	return base.Add(quotient * dsec * time.Second)
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

// CalculateDaysWithBaseline 計算兩個時間相差幾天, 但以 base 指定的時間為日界
//
// base 應該是個介於 0 ~ 23:59:59 的時間長度, 例如 base 為 5 小時的話, 表示以凌晨 5 點為日界
func CalculateDaysWithBaseline(t1, t2 time.Time, base time.Duration) int {
	return CalculateDays(t1.Add(-base), t2.Add(-base))
}

// clampHour 將輸入的小時正規化到合法範圍 [0, 23]
func clampHour(hour int) int {
	return min(max(hour, 0), 23) //nolint: mnd
}

// clampWeek 將輸入的星期正規化到合法範圍 [0, 6]
func clampWeek(wday int) int {
	return ((wday % 7) + 7) % 7 //nolint: mnd
}

// clampMDay 將輸入的日期正規化到合法範圍 [1, 31]
func clampMDay(mday int) int {
	return min(max(mday, 1), 31) //nolint: mnd
}

// clampMonth 將輸入的月份正規化到合法範圍 [1, 12]
func clampMonth(month int) int {
	return min(max(month, 1), 12) //nolint: mnd
}

var zone *time.Location // 時區設定, 若未設定則預設為 UTC, 所有透過本套件建構/解析的時間, 皆會套用此時區
