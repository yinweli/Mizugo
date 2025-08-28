package cryptos

import (
	"bytes"
	"fmt"

	"github.com/yinweli/Mizugo/v2/mizugos/helps"
)

const (
	DesKeySize           = 8                    // DES 密鑰長度
	DesKeyLetter         = helps.StrNumberAlpha // 產生隨機密鑰時可使用的字元集
	PaddingZero  Padding = iota                 // Zero padding
	PaddingPKCS7                                // PKCS7 padding
)

// Padding 填充模式
type Padding = int

// RandDesKey 產生隨機 DES 密鑰
func RandDesKey() []byte {
	return []byte(helps.RandString(DesKeySize, DesKeyLetter))
}

// RandDesKeyString 產生隨機 DES 密鑰字串
func RandDesKeyString() string {
	return string(RandDesKey())
}

// pad 對資料進行填充, 使其長度符合區塊大小的倍數
func pad(padding Padding, source []byte, blockSize int) (result []byte, err error) {
	switch padding {
	case PaddingZero:
		if result, err = zeroPad(source, blockSize); err != nil {
			return nil, fmt.Errorf("pad: %w", err)
		} // if

		return result, nil

	case PaddingPKCS7:
		if result, err = pkcs7Pad(source, blockSize); err != nil {
			return nil, fmt.Errorf("pad: %w", err)
		} // if

		return result, nil

	default:
		return nil, fmt.Errorf("pad: unknown padding %v", padding)
	} // switch
}

// unpad 對資料進行反填充, 移除多餘的補齊資料
func unpad(padding Padding, source []byte) (result []byte, err error) {
	switch padding {
	case PaddingZero:
		if result, err = zeroUnpad(source); err != nil {
			return nil, fmt.Errorf("unpad: %w", err)
		} // if

		return result, nil

	case PaddingPKCS7:
		if result, err = pkcs7Unpad(source); err != nil {
			return nil, fmt.Errorf("unpad: %w", err)
		} // if

		return result, nil

	default:
		return nil, fmt.Errorf("unpad: unknown padding %v", padding)
	} // switch
}

// zeroPad 使用 Zero padding 將資料補齊; 不足區塊大小的部分以 0 填滿
func zeroPad(source []byte, blockSize int) (result []byte, err error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("zeroPad: block size <= 0")
	} // if

	size := blockSize - len(source)%blockSize
	text := bytes.Repeat([]byte{0}, size)
	source = append(source, text...)
	return source, nil
}

// zeroUnpad 移除 Zero padding 的補齊資料
func zeroUnpad(source []byte) (result []byte, err error) { //nolint:unparam
	i := len(source)

	for i > 0 && source[i-1] == 0 {
		i--
	} // for

	return source[:i], nil
}

// pkcs7Pad 使用 PKCS7 padding 將資料補齊; 不足區塊大小的部分以「補齊長度」作為填充值
func pkcs7Pad(source []byte, blockSize int) (result []byte, err error) {
	if blockSize <= 0 {
		return nil, fmt.Errorf("pkcs7Pad: block size <= 0")
	} // if

	size := blockSize - len(source)%blockSize
	text := bytes.Repeat([]byte{byte(size)}, size)
	source = append(source, text...)
	return source, nil
}

// pkcs7Unpad 移除 PKCS7 padding 的補齊資料; 會檢查補齊長度是否合法, 並透過聚合比對避免時序側信道攻擊
func pkcs7Unpad(source []byte) (result []byte, err error) {
	length := len(source)

	if length == 0 {
		return nil, fmt.Errorf("pkcs7Unpad: source empty")
	} // if

	size := int(source[length-1])

	if size == 0 || size > length {
		return nil, fmt.Errorf("pkcs7Unpad: invalid size")
	} // if

	// 檢查填充數值是否正確, 並使用聚合比對, 避免時序側信道風險
	check := 0

	for i := 0; i < size; i++ {
		check |= int(source[length-1-i] ^ byte(size))
	} // for

	if check != 0 {
		return nil, fmt.Errorf("pkcs7Unpad: invalid padding")
	} // if

	return source[:length-size], nil
}
