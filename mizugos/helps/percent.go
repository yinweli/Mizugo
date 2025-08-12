package helps

const (
	PercentRatio100 = 100   // 百分比例值
	PercentRatio1K  = 1000  // 千分比例值
	PercentRatio10K = 10000 // 萬分比例值
)

// NewP100 建立百分比計算器
func NewP100() *Percent {
	return NewPercent(PercentRatio100)
}

// NewP1K 建立千分比計算器
func NewP1K() *Percent {
	return NewPercent(PercentRatio1K)
}

// NewP10K 建立萬分比計算器
func NewP10K() *Percent {
	return NewPercent(PercentRatio10K)
}

// NewPercent 建立比例計算器
func NewPercent(base int32) *Percent {
	return &Percent{
		base: base,
	}
}

// Percent 比例計算器
type Percent struct {
	base int32 // 基準值
	per  int32 // 比例值
}

// Rounder 進位函式類型
type Rounder func(value float64) float64

// Base 取得基準值
func (this *Percent) Base() int32 {
	return this.base
}

// Set 設定比例值
func (this *Percent) Set(per int32) *Percent {
	this.per = per
	return this
}

// SetBase 設定比例值為基準值
func (this *Percent) SetBase() *Percent {
	this.Set(this.base)
	return this
}

// Add 增加比例值
func (this *Percent) Add(per int32) *Percent {
	this.per += per
	return this
}

// Del 減少比例值
func (this *Percent) Del(per int32) *Percent {
	this.per -= per
	return this
}

// Mul 乘以比例值
func (this *Percent) Mul(per int32) *Percent {
	this.per *= per
	return this
}

// Div 除以比例值
func (this *Percent) Div(per int32) *Percent {
	if per != 0 {
		this.per /= per
	} // if

	return this
}

// Get 取得比例值
func (this *Percent) Get() int32 {
	return this.per
}

// Calc 計算結果, input為輸入值, round為使用哪個函式來計算進位(math.Round / math.Ceil / math.Floor)
func (this *Percent) Calc(input int32, round Rounder) int {
	return int(this.calc(float64(input), round))
}

// Calc32 計算結果, input為輸入值, round為使用哪個函式來計算進位(math.Round / math.Ceil / math.Floor)
func (this *Percent) Calc32(input int32, round Rounder) int32 {
	return int32(this.calc(float64(input), round))
}

// Calc64 計算結果, input為輸入值, round為使用哪個函式來計算進位(math.Round / math.Ceil / math.Floor)
func (this *Percent) Calc64(input int64, round Rounder) int64 {
	return int64(this.calc(float64(input), round))
}

// calc 計算結果
func (this *Percent) calc(input float64, round Rounder) float64 {
	if this.base == 0 {
		panic("Percent.calc: base is zero")
	} // if

	if round == nil {
		panic("Percent.calc: round is nil")
	} // if

	return round(input * float64(this.per) / float64(this.base))
}
