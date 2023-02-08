package cryptos

import (
	"crypto/cipher"
	"crypto/des"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

const DesKeySize = 8 // des密鑰長度

// DesECBEncrypt des-ecb加密, 注意key的長度必須是 DesKeySize
func DesECBEncrypt(padding Padding, key, src []byte) (out []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb encrypt: key len must %v", DesKeySize)
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb encrypt: %w", err)
	} // if

	blockSize := block.BlockSize()
	src = pad(padding, src, blockSize)
	out = make([]byte, len(src))
	dst := out

	for len(src) > 0 {
		block.Encrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	return out, nil
}

// DesECBDecrypt des-ecb解密, 注意key的長度必須是 DesKeySize
func DesECBDecrypt(padding Padding, key, src []byte) (out []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb decrypt: key len must %v", DesKeySize)
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb decrypt: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("des-ecb decrypt: src not full blocks")
	} // if

	out = make([]byte, len(src))
	dst := out

	for len(src) > 0 {
		block.Decrypt(dst, src[:blockSize])
		src = src[blockSize:]
		dst = dst[blockSize:]
	} // for

	out = unpad(padding, out)
	return out, nil
}

// DesCBCEncrypt des-cbc加密, 注意key的長度必須是 DesKeySize
func DesCBCEncrypt(padding Padding, key, iv, src []byte) (out []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-cbc encrypt: key len must %v", DesKeySize)
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-cbc encrypt: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(iv) != blockSize {
		return nil, fmt.Errorf("des-cbc encrypt: iv len must %v", blockSize)
	} // if

	src = pad(padding, src, blockSize)
	out = make([]byte, len(src))
	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(out, src)
	return out, nil
}

// DesCBCDecrypt des-cbc解密, 注意key的長度必須是 DesKeySize
func DesCBCDecrypt(padding Padding, key, iv, src []byte) (out []byte, err error) {
	if len(key) != DesKeySize {
		return nil, fmt.Errorf("des-cbc decrypt: key len must %v", DesKeySize)
	} // if

	block, err := des.NewCipher(key)

	if err != nil {
		return nil, fmt.Errorf("des-cbc decrypt: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(iv) != blockSize {
		return nil, fmt.Errorf("des-cbc decrypt: iv len must %v", blockSize)
	} // if

	if len(src)%blockSize != 0 {
		return nil, fmt.Errorf("des-cbc decrypt: src not full blocks")
	} // if

	out = make([]byte, len(src))
	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypter.CryptBlocks(out, src)
	out = unpad(padding, out)
	return out, nil
}

// RandDesKey 產生隨機des密鑰
func RandDesKey() []byte {
	return []byte(utils.RandString(DesKeySize))
}

// RandDesKeyString 產生隨機des密鑰字串
func RandDesKeyString() string {
	return string(RandDesKey())
}
