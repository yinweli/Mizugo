package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

// MD5String 取得MD5字串
func MD5String(input []byte) string {
	result := md5.Sum(input)
	return hex.EncodeToString(result[:])
}

// Base64Encode base64加密
func Base64Encode(input []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(input))
}

// Base64Decode base64解密
func Base64Decode(input []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(input))
}
