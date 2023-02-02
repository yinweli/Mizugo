package cryptos

import (
	"crypto/des"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

const DesKeySize = 8 // Des密鑰長度

// DesECBEncrypt Des-ECB加密, 注意key只能是8位陣列
func DesECBEncrypt(padding Padding, key, input []byte) (result []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb encrypt: key len must 8")
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb encrypt: %w", err)
	} // if

	blockSize := block.BlockSize()
	src := pad(padding, input, blockSize)
	out := make([]byte, len(src))
	dst := out

	for len(src) > 0 {
		block.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	return out, nil
}

// DesECBDecrypt Des-ECB解密, 注意key只能是8位陣列
func DesECBDecrypt(padding Padding, key, input []byte) (result []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb decrypt: key len must 8")
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb decrypt: %w", err)
	} // if

	blockSize := block.BlockSize()
	src := input
	out := make([]byte, len(src))
	dst := out

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("des-ecb decrypt: input not full blocks")
	} // if

	for len(src) > 0 {
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	out = unpad(padding, out)
	return out, nil
}

// TODO: Des-CBC
// TODO: C#的Des-EBC & Des-CBC
// https://blog.clarence.tw/2020/12/28/golang_implements_aes_ecb_and_pkcs7_pkgs5/
// https://www.cnblogs.com/shanfeng1000/p/14808574.html
// https://www.itread01.com/articles/1475821568.html
// https://ithelp.ithome.com.tw/articles/10250386

// RandDesKey 產生隨機Des密鑰
func RandDesKey() []byte {
	return []byte(utils.RandString(DesKeySize))
}

// RandDesKeyString 產生隨機Des密鑰字串
func RandDesKeyString() string {
	return string(RandDesKey())
}
