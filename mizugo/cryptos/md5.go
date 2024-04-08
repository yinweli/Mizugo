package cryptos

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5String 取得md5字串
func MD5String(input []byte) string {
	result := md5.Sum(input)
	return hex.EncodeToString(result[:])
}
