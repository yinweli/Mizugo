package helps

import (
	"time"
)

// NewCurrent 建立目前時間資料
func NewCurrent() *Current {
	now := Time()
	return &Current{
		base:  now,
		begin: now,
	}
}

// Current 目前時間資料, 用在目標資料中,
// 當目標的內部函式要用到時間相關的功能時, 需要通過 Current 提供的功能來取得,
// 這樣就能讓單元測試時可以通過改變基準時間, 來改變時間
type Current struct {
	base  time.Time // 基準時間
	begin time.Time // 起點時間
}

// Curr 取得自己物件
func (this *Current) Curr() *Current {
	return this
}

// SetBaseNow 設定基準時間
func (this *Current) SetBaseNow() {
	this.base = Time()
}

// SetBaseTime 設定基準時間
func (this *Current) SetBaseTime(now time.Time) {
	this.base = now
}

// SetBaseDate 設定基準時間
func (this *Current) SetBaseDate(year int, month time.Month, day, hour, min, sec int) {
	this.base = Date(year, month, day, hour, min, sec, 0)
}

// SetBaseDay 設定基準時間
func (this *Current) SetBaseDay(year int, month time.Month, day int) {
	this.base = Date(year, month, day, 0, 0, 0, 0)
}

// AddBaseTime 增加基準時間
func (this *Current) AddBaseTime(duration time.Duration) {
	this.base = this.base.Add(duration)
}

// GetBase 取得基準時間
func (this *Current) GetBase() time.Time {
	return this.base
}

// GetTime 取得時間
func (this *Current) GetTime() time.Time {
	return this.base.Add(Time().Sub(this.begin))
}
