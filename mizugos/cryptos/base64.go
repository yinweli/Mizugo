package cryptos

import (
	"encoding/base64"
)

// Base64Encode base64加密
func Base64Encode(input []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(input))
}

// Base64Decode base64解密
func Base64Decode(input []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(input))
}
