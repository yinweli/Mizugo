package cryptos

import (
	"crypto/des"
	"fmt"
)

// NewDesECB 建立 DES-ECB 編碼/解碼器
//
// 參數:
//   - padding: 填充模式(Zero padding 或 PKCS7 padding)
//   - key: 密鑰字串, 長度必須等於 DesKeySize
func NewDesECB(padding Padding, key string) *DesECB {
	return &DesECB{
		padding: padding,
		key:     []byte(key),
	}
}

// DesECB DES-ECB 編碼/解碼器
//
// 注意: ECB 模式無法隱藏明文模式, 通常不建議用於安全性要求高的場合
type DesECB struct {
	padding Padding // 填充模式
	key     []byte  // 密鑰, 長度必須是 DesKeySize
}

// Encode 編碼
func (this *DesECB) Encode(input any) (output any, err error) {
	if len(this.key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb encode: key len must %v", DesKeySize)
	} // if

	if input == nil {
		return nil, fmt.Errorf("des-ecb encode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("des-ecb encode: input type")
	} // if

	if len(source) == 0 {
		return nil, fmt.Errorf("des-ecb encode: input empty")
	} // if

	block, err := des.NewCipher(this.key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb encode: %w", err)
	} // if

	blockSize := block.BlockSize()

	if source, err = pad(this.padding, source, blockSize); err != nil {
		return nil, fmt.Errorf("des-ecb encode: %w", err)
	} // if

	result := make([]byte, len(source))
	dst := result

	for len(source) > 0 {
		block.Encrypt(dst, source[:blockSize])
		source = source[blockSize:]
		dst = dst[blockSize:]
	} // for

	return result, nil
}

// Decode 解碼
func (this *DesECB) Decode(input any) (output any, err error) {
	if len(this.key) != DesKeySize {
		return nil, fmt.Errorf("des-ecb decode: key len must %v", DesKeySize)
	} // if

	if input == nil {
		return nil, fmt.Errorf("des-ecb decode: input nil")
	} // if

	source, ok := input.([]byte)

	if ok == false {
		return nil, fmt.Errorf("des-ecb decode: input type")
	} // if

	if len(source) == 0 {
		return nil, fmt.Errorf("des-ecb decode: input empty")
	} // if

	block, err := des.NewCipher(this.key)

	if err != nil {
		return nil, fmt.Errorf("des-ecb decode: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(source)%blockSize != 0 {
		return nil, fmt.Errorf("des-ecb decode: src not full blocks")
	} // if

	result := make([]byte, len(source))
	dst := result

	for len(source) > 0 {
		block.Decrypt(dst, source[:blockSize])
		source = source[blockSize:]
		dst = dst[blockSize:]
	} // for

	if result, err = unpad(this.padding, result); err != nil {
		return nil, fmt.Errorf("des-ecb decode: %w", err)
	} // if

	return result, nil
}
