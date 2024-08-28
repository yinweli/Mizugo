package cryptos

import (
	"crypto/des"
	"fmt"
)

// NewDesECB 建立des-ecb加密/解密器
func NewDesECB(padding Padding, key string) *DesECB {
	return &DesECB{
		padding: padding,
		key:     []byte(key),
	}
}

// DesECB des-ecb加密/解密器
type DesECB struct {
	padding Padding // 填充模式
	key     []byte  // 密鑰, 長度必須是 DesKeySize
}

// Encode 加密
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
	source = pad(this.padding, source, blockSize)
	target := make([]byte, len(source))
	dst := target

	for len(source) > 0 {
		block.Encrypt(dst, source[:blockSize])
		source = source[blockSize:]
		dst = dst[blockSize:]
	} // for

	return target, nil
}

// Decode 解密
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

	target := make([]byte, len(source))
	dst := target

	for len(source) > 0 {
		block.Decrypt(dst, source[:blockSize])
		source = source[blockSize:]
		dst = dst[blockSize:]
	} // for

	target = unpad(this.padding, target)
	return target, nil
}
