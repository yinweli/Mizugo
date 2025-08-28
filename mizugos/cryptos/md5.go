package cryptos

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5String 將輸入的陣列計算 MD5 雜湊值, 並回傳其 16 進位字串
func MD5String(input []byte) string {
	result := md5.Sum(input)
	return hex.EncodeToString(result[:])
}
