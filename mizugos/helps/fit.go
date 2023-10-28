package helps

// NewFit 建立數值檢測資料
func NewFit[T int | int32 | int64](maxCap, minCap FitCap[T]) *Fit[T] {
	return &Fit[T]{
		maxCap: maxCap,
		minCap: minCap,
	}
}

// Fit 數值檢測資料
type Fit[T int | int32 | int64] struct {
	maxCap FitCap[T] // 取得最大值函式
	minCap FitCap[T] // 取得最小值函式
}

// Check 數值檢測
//   - 當輸入值在範圍時: result = 輸入值, remain = 0
//   - 當輸入值大於上限時: result = 上限值, remain = 溢出值(正值)
//   - 當輸入值小於下限時: result = 下限值, remain = 溢出值(負值)
func (this *Fit[T]) Check(source T, modify ...T) (result, remain, added T) {
	fin := int64(source)

	for _, itor := range modify {
		fin += int64(itor)
	} // for

	valueMax := int64(this.maxCap.do())
	valueMin := int64(this.minCap.do())

	if fin > valueMax {
		return T(valueMax), T(fin - valueMax), T(valueMax) - source
	} // if

	if fin < valueMin {
		return T(valueMin), T(fin - valueMin), T(valueMin) - source
	} // if

	return T(fin), 0, T(fin) - source
}

// FitCap 數值限制函式類型
type FitCap[T int | int32 | int64] func() T

// do 執行數值限制函式
func (this FitCap[T]) do() T {
	if this != nil {
		return this()
	} // if

	return 0
}
