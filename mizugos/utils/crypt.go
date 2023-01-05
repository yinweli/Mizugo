package utils

import (
	"bytes"
	"crypto/des"
	"fmt"
)

const DesKeySize = 8 // Des密鑰長度

// DesEncrypt Des加密處理, 注意key只能是8位陣列
// 這裡選用 https://blog.csdn.net/wade3015/article/details/84454836 提供的方式, 數據填充方式用zeropad
func DesEncrypt(key, input []byte) (result []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des encrypt: key len must 8")
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des encrypt: %w", err)
	} // if

	blockSize := block.BlockSize()
	src := input
	padSize := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{0}, padSize)
	src = append(src, padText...)

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("des encrypt: need a multiple of the blocksize")
	} // if

	out := make([]byte, len(src))
	dst := out

	for len(src) > 0 {
		block.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	return out, nil
}

// DesDecrypt Des解密處理, 注意key只能是8位陣列
// 這裡選用 https://blog.csdn.net/wade3015/article/details/84454836 提供的方式, 數據填充方式用zeropad
func DesDecrypt(key, input []byte) (result []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des decrypt: key len must 8")
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des decrypt: %w", err)
	} // if

	blockSize := block.BlockSize()
	src := input
	out := make([]byte, len(src))
	dst := out

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("des decrypt: input not full blocks")
	} // if

	for len(src) > 0 {
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	out = bytes.TrimFunc(out, func(r rune) bool {
		return r == rune(0)
	})
	return out, nil
}

// DesKeyRand 產生隨機Des密鑰
func DesKeyRand() []byte {
	return []byte(RandString(DesKeySize))
}
