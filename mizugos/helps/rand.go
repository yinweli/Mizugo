package helps

import (
	"bytes"
	randr "crypto/rand"
	"math"
	"math/big"
	"math/rand/v2"
)

// RandInt 取得偽隨機數值, 不會有負值
func RandInt() int {
	return rand.Int() //nolint:gosec
}

// RandIntn 取得範圍內偽隨機數值
func RandIntn(minimum, maximum int) int {
	m := max(minimum, maximum)
	n := min(minimum, maximum)
	return rand.N(m-n+1) + n //nolint:gosec
}

// RandInt32 取得偽隨機數值, 不會有負值
func RandInt32() int32 {
	return rand.Int32() //nolint:gosec
}

// RandInt32n 取得範圍內偽隨機數值
func RandInt32n(minimum, maximum int32) int32 {
	m := max(minimum, maximum)
	n := min(minimum, maximum)
	return rand.N(m-n+1) + n //nolint:gosec
}

// RandInt64 取得偽隨機數值, 不會有負值
func RandInt64() int64 {
	return rand.Int64() //nolint:gosec
}

// RandInt64n 取得範圍內偽隨機數值
func RandInt64n(minimum, maximum int64) int64 {
	m := max(minimum, maximum)
	n := min(minimum, maximum)
	return rand.N(m-n+1) + n //nolint:gosec
}

// RandReal64 取得真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64() int64 {
	value, _ := randr.Int(randr.Reader, big.NewInt(math.MaxInt64))
	return value.Int64()
}

// RandReal64n 取得範圍內真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64n(minimum, maximum int64) int64 {
	mini := min(minimum, maximum)
	maxi := max(minimum, maximum)
	value, _ := randr.Int(randr.Reader, big.NewInt(maxi-mini+1))
	return mini + value.Int64()
}

// RandString 取得指定位數的隨機字串, 從輸入的 letter 集合中隨機挑選字元
func RandString(length int, letter string) string {
	builder := bytes.Buffer{}
	lettern := len(letter)

	for i := 0; i < length; i++ {
		index := rand.N(lettern) //nolint:gosec
		builder.WriteByte(letter[index])
	} // for

	return builder.String()
}

// RandStringDefault 取得隨機字串, 長度為 10, 字元集為 StrNumberAlpha
func RandStringDefault() string {
	const length = 10
	const letter = StrNumberAlpha
	return RandString(length, letter)
}
