package helps

// NewFit 建立數值檢測資料
func NewFit[T ~int | ~int32 | ~int64](maximum, minimum func() T) *Fit[T] {
	return &Fit[T]{
		maximum: maximum,
		minimum: minimum,
	}
}

// Fit 數值檢測資料
//
// 適用於有「上下限」概念的數值(如血量、能量、經驗值), 透過 maximum / minimum 函式動態提供上下界
//
// 當沒有提供 maximum / minimum 時, 會以 0 為預設值
type Fit[T ~int | ~int32 | ~int64] struct {
	maximum func() T // 取得最大值函式
	minimum func() T // 取得最小值函式
}

// Check 檢測數值範圍, 並計算調整後的結果
//
// 計算方式:
//   - 輸入值 = source + sum(modify...)
//   - 若輸入值在 [minimum, maximum] 之間:
//     result = 輸入值, remain = 0, added = 輸入值 - source
//   - 若輸入值 > maximum:
//     result = maximum, remain = 輸入值 - maximum (正值), added = maximum - source
//   - 若輸入值 < minimum:
//     result = minimum, remain = 輸入值 - minimum (負值), added = minimum - source
//
// 回傳值:
//   - result: 實際落在範圍內的值
//   - remain: 超出範圍的溢出量(大於上限時為正, 小於下限時為負)
//   - added : 相對於原始 source 的淨增減量(即 result - source)
func (this *Fit[T]) Check(source T, modify ...T) (result, remain, added T) {
	fin := int64(source)

	for _, itor := range modify {
		fin += int64(itor)
	} // for

	valueMax := int64(0)

	if this.maximum != nil {
		valueMax = int64(this.maximum())
	} // if

	valueMin := int64(0)

	if this.minimum != nil {
		valueMin = int64(this.minimum())
	} // if

	if fin > valueMax {
		return T(valueMax), T(fin - valueMax), T(valueMax) - source
	} // if

	if fin < valueMin {
		return T(valueMin), T(fin - valueMin), T(valueMin) - source
	} // if

	return T(fin), 0, T(fin) - source
}
