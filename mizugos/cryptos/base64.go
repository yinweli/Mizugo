package cryptos

import (
	"encoding/base64"
	"fmt"
)

// NewBase64 建立 Base64 編碼/解碼器
func NewBase64() *Base64 {
	return &Base64{}
}

// Base64 編碼/解碼器
type Base64 struct {
}

// Encode 編碼
func (this *Base64) Encode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("base64 encode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("base64 encode: input type")
	} // if

	return []byte(base64.StdEncoding.EncodeToString(source)), nil
}

// Decode 解碼
func (this *Base64) Decode(input any) (output any, err error) {
	if input == nil {
		return nil, fmt.Errorf("base64 decode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("base64 decode: input type")
	} // if

	return base64.StdEncoding.DecodeString(string(source))
}
