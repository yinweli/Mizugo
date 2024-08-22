package cryptos

import (
	"crypto/cipher"
	"crypto/des"
	"fmt"
)

// NewDesCBC 建立des-cbc加密/解密器
func NewDesCBC(padding Padding, key, iv string) *DesCBC {
	return &DesCBC{
		padding: padding,
		key:     []byte(key),
		iv:      []byte(iv),
	}
}

// DesCBC des-cbc加密/解密器
type DesCBC struct {
	padding Padding // 填充模式
	key     []byte  // 密鑰, 長度必須是 DesKeySize
	iv      []byte  // 向量值
}

// Encode 加密
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

	block, err := des.NewCipher(this.key) //nolint:gosec

	if err != nil {
		return nil, fmt.Errorf("des-cbc encode: %w", err)
	} // if

	blockSize := block.BlockSize()

	if len(this.iv) != blockSize {
		return nil, fmt.Errorf("des-cbc encode: iv len must %v", blockSize)
	} // if

	source = pad(this.padding, source, blockSize)
	target := make([]byte, len(source))
	encrypter := cipher.NewCBCEncrypter(block, this.iv)
	encrypter.CryptBlocks(target, source)
	return target, nil
}

// Decode 解密
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

	block, err := des.NewCipher(this.key) //nolint:gosec

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

	target := make([]byte, len(source))
	decrypter := cipher.NewCBCDecrypter(block, this.iv)
	decrypter.CryptBlocks(target, source)
	target = unpad(this.padding, target)
	return target, nil
}
