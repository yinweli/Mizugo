package helps

import (
	"bytes"
	randr "crypto/rand"
	"math"
	"math/big"
	"sync"
	"time"

	"golang.org/x/exp/rand"
)

// RandSeed 設定偽隨機種子
func RandSeed(seed int64) {
	rnLock.Lock()
	defer rnLock.Unlock()
	rn.Seed(uint64(seed))
}

// RandSeedTime 設定偽隨機種子, 以目前的時間值來設定
func RandSeedTime() {
	rnLock.Lock()
	defer rnLock.Unlock()
	rn.Seed(uint64(time.Now().UnixNano()))
}

// RandSource 取得隨機來源
func RandSource() rand.Source {
	return rns
}

// RandInt 取得偽隨機數值, 不會有負值
func RandInt() int {
	rnLock.Lock()
	defer rnLock.Unlock()
	return rn.Int()
}

// RandIntn 取得範圍內偽隨機數值
func RandIntn(minimum, maximum int) int {
	rnLock.Lock()
	defer rnLock.Unlock()
	return minimum + rn.Intn(maximum-minimum+1)
}

// RandInt32 取得偽隨機數值, 不會有負值
func RandInt32() int32 {
	rnLock.Lock()
	defer rnLock.Unlock()
	return rn.Int31()
}

// RandInt32n 取得範圍內偽隨機數值
func RandInt32n(minimum, maximum int32) int32 {
	rnLock.Lock()
	defer rnLock.Unlock()
	return minimum + rn.Int31n(maximum-minimum+1)
}

// RandInt64 取得偽隨機數值, 不會有負值
func RandInt64() int64 {
	rnLock.Lock()
	defer rnLock.Unlock()
	return rn.Int63()
}

// RandInt64n 取得範圍內偽隨機數值
func RandInt64n(minimum, maximum int64) int64 {
	rnLock.Lock()
	defer rnLock.Unlock()
	return minimum + rn.Int63n(maximum-minimum+1)
}

// RandReal64 取得真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64() int64 {
	value, _ := randr.Int(randr.Reader, big.NewInt(math.MaxInt64))
	return value.Int64()
}

// RandReal64n 取得範圍內真隨機數值, 不需要事先設定隨機種子, 但是速度大約比偽隨機慢10倍
func RandReal64n(minimum, maximum int64) int64 {
	value, _ := randr.Int(randr.Reader, big.NewInt(maximum-minimum+1))
	return minimum + value.Int64()
}

// RandString 取得指定位數的隨機字串
func RandString(length int, letter string) string {
	builder := bytes.Buffer{}
	lettern := len(letter)

	rnLock.Lock()
	defer rnLock.Unlock()

	for i := 0; i < length; i++ {
		index := rn.Intn(lettern)
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

var rn = rand.New(rns)                                  // 隨機產生器
var rns = rand.NewSource(uint64(time.Now().UnixNano())) // 隨機來源
var rnLock sync.Mutex                                   // 執行緒鎖
