package utils

import (
	"bytes"
	"crypto/rand"
	"math/big"
)

// RandString 取得指定位數的隨機字串
func RandString(count int) string {
	const letter = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	builder := bytes.Buffer{}
	max := big.NewInt(int64(len(letter)))

	for i := 0; i < count; i++ {
		index, _ := rand.Int(rand.Reader, max)
		builder.WriteByte(letter[int(index.Int64())])
	} // for

	return builder.String()
}
