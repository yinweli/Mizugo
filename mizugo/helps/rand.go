package helps

import (
	"bytes"
	randr "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"time"
)

// RandSeed 設定偽隨機種子
func RandSeed(seed int64) {
	rand.Seed(seed)
}

// RandSeedTime 設定偽隨機種子, 以目前的時間值來設定
func RandSeedTime() {
	rand.Seed(time.Now().UnixNano())
}

// RandInt 取得偽隨機數值, 不會有負值
func RandInt() int {
	return rand.Int() //nolint:gosec
}

// RandIntn 取得範圍內偽隨機數值
func RandIntn(min, max int) int {
	return min + rand.Intn(max-min+1) //nolint:gosec
}

// RandInt32 取得偽隨機數值, 不會有負值
func RandInt32() int32 {
	return rand.Int31() //nolint:gosec
}

// RandInt32n 取得範圍內偽隨機數值
func RandInt32n(min, max int32) int32 {
	return min + rand.Int31n(max-min+1) //nolint:gosec
}

// RandInt64 取得偽隨機數值, 不會有負值
func RandInt64() int64 {
	return rand.Int63() //nolint:gosec
}

// RandInt64n 取得範圍內偽隨機數值
func RandInt64n(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) //nolint:gosec
}

// RandReal64 取得真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64() int64 {
	value, _ := randr.Int(randr.Reader, big.NewInt(math.MaxInt64))
	return value.Int64()
}

// RandReal64n 取得範圍內真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64n(min, max int64) int64 {
	value, _ := randr.Int(randr.Reader, big.NewInt(max-min+1))
	return min + value.Int64()
}

// RandString 取得指定位數的隨機字串
func RandString(length int, letter string) string {
	builder := bytes.Buffer{}
	lettern := len(letter) - 1

	for i := 0; i < length; i++ {
		index := RandIntn(0, lettern)
		builder.WriteByte(letter[index])
	} // for

	return builder.String()
}

// RandStringDefault 取得預設配置的隨機字串
func RandStringDefault() string {
	const length = 10
	const letter = StrNumberAlpha
	return RandString(length, letter)
}
