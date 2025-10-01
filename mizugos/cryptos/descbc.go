package cryptos

import (
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

// NewDesCBC 建立 DES-CBC 編碼/解碼器
//
// 參數:
//   - padding: 填充模式(Zero padding 或 PKCS7 padding)
//   - key: 密鑰字串, 長度必須等於 DesKeySize
//   - iv: 初始化向量, 長度必須等於 DesKeySize
func NewDesCBC(padding Padding, key, iv string) *DesCBC {
	return &DesCBC{
		padding: padding,
		key:     []byte(key),
		iv:      []byte(iv),
	}
}

// DesCBC DES-CBC 編碼/解碼器
type DesCBC struct {
	padding Padding // 填充模式
	key     []byte  // 密鑰, 長度必須是 DesKeySize
	iv      []byte  // 初始化向量
}

// Encode 編碼
func (this *DesCBC) Encode(input any) (output any, err error) {
	if len(this.key) != DesKeySize {
		return nil, fmt.Errorf("des-cbc encode: key len must %v", DesKeySize)
	} // if

	if len(this.iv) == 0 {
		return nil, fmt.Errorf("des-cbc encode: iv nil")
	} // if

	if input == nil {
		return nil, fmt.Errorf("des-cbc encode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("des-cbc encode: input type")
	} // if

	if len(source) == 0 {
		return nil, fmt.Errorf("des-cbc encode: input empty")
	} // if

	block, err := des.NewCipher(this.key)

	if err != nil {
		return nil, fmt.Errorf("des-cbc encode: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(this.iv) != blockSize {
		return nil, fmt.Errorf("des-cbc encode: iv len must %v", blockSize)
	} // if

	if source, err = pad(this.padding, source, blockSize); err != nil {
		return nil, fmt.Errorf("des-cbc encode: %w", err)
	} // if

	result := make([]byte, len(source))
	encrypter := cipher.NewCBCEncrypter(block, this.iv)
	encrypter.CryptBlocks(result, source)
	return result, nil
}

// Decode 解碼
func (this *DesCBC) Decode(input any) (output any, err error) {
	if len(this.key) != DesKeySize {
		return nil, fmt.Errorf("des-cbc decode: key len must %v", DesKeySize)
	} // if

	if len(this.iv) == 0 {
		return nil, fmt.Errorf("des-cbc decode: iv nil")
	} // if

	if input == nil {
		return nil, fmt.Errorf("des-cbc decode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("des-cbc decode: input type")
	} // if

	if len(source) == 0 {
		return nil, fmt.Errorf("des-cbc decode: input empty")
	} // if

	block, err := des.NewCipher(this.key)

	if err != nil {
		return nil, fmt.Errorf("des-cbc decode: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(this.iv) != blockSize {
		return nil, fmt.Errorf("des-cbc decode: iv len must %v", blockSize)
	} // if

	if len(source)%blockSize != 0 {
		return nil, fmt.Errorf("des-cbc decode: src not full blocks")
	} // if

	result := make([]byte, len(source))
	decrypter := cipher.NewCBCDecrypter(block, this.iv)
	decrypter.CryptBlocks(result, source)

	if result, err = unpad(this.padding, result); err != nil {
		return nil, fmt.Errorf("des-cbc decode: %w", err)
	} // if

	return result, nil
}
