package cryptos

import (
	"bytes"
)

const ( // 填充模式編號
	PaddingZero  Padding = iota // zeropad填充
	PaddingPKCS7                // pkcs7填充
)

// Padding 填充模式編號
type Padding = int

// pad 填充
func pad(padding Padding, source []byte, blockSize int) []byte {
	switch padding {
	case PaddingZero:
		return padZero(source, blockSize)

	case PaddingPKCS7:
		return padPKCS7(source, blockSize)

	default:
		return nil
	} // switch
}

// unpad 反填充
func unpad(padding Padding, source []byte) []byte {
	switch padding {
	case PaddingZero:
		return unpadZero(source)

	case PaddingPKCS7:
		return unpadPKCS7(source)

	default:
		return nil
	} // switch
}

// padZero zeropad填充
func padZero(source []byte, blockSize int) []byte {
	size := blockSize - len(source)%blockSize
	text := bytes.Repeat([]byte{0}, size)
	return append(source, text...)
}

// unpadZero zeropad反填充
func unpadZero(source []byte) []byte {
	return bytes.TrimFunc(source, func(r rune) bool {
		return r == rune(0)
	})
}

// padPKCS7 pkcs7填充
func padPKCS7(source []byte, blockSize int) []byte {
	size := blockSize - len(source)%blockSize
	text := bytes.Repeat([]byte{byte(size)}, size)
	return append(source, text...)
}

// unpadPKCS7 pkcs7反填充
func unpadPKCS7(source []byte) []byte {
	length := len(source)
	size := source[length-1]
	return source[:length-int(size)]
}
