package helps

import (
	"fmt"
	"strings"
)

const (
	// Base58 使用的字元表, 排除了容易混淆的字母 o O l I
	Base58 = "0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
	// Base80 使用的字元表, 涵蓋數字, 大小寫字母與常見符號
	Base80 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ._~+-*=?!|@#$%^&<>"
)

// ToBase58 將 uint64 轉為 58 進制字串
func ToBase58(input uint64) string {
	return ToBaseN(Base58, input)
}

// FromBase58 將 58 進制字串轉為 uint64
func FromBase58(input string) (result uint64, err error) {
	result, err = FromBaseN(Base58, input)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// ToBase80 將 uint64 轉為 80 進制字串
func ToBase80(input uint64) string {
	return ToBaseN(Base80, input)
}

// FromBase80 將 80 進制字串轉為 uint64
func FromBase80(input string) (result uint64, err error) {
	result, err = FromBaseN(Base80, input)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// ToBaseN 將 uint64 轉為自訂進制字串, model 為字元表, 字元數量決定進制的基數, 字元表中不能有重複字元
func ToBaseN(model string, input uint64) string {
	if input == 0 {
		return "0"
	} // if

	base := uint64(len(model))
	encoded := []byte{}

	for input > 0 {
		mod := input % base
		encoded = append(encoded, model[mod])
		input /= base
	} // for

	for i, j := 0, len(encoded)-1; i < j; i, j = i+1, j-1 {
		encoded[i], encoded[j] = encoded[j], encoded[i]
	} // for

	return string(encoded)
}

// FromBaseN 將自訂進制字串轉換回 uint64, model 為字元表, 必須與 ToBaseN 使用的相同, 字元表中不能有重複字元
func FromBaseN(model, input string) (result uint64, err error) {
	if input == "" {
		return 0, fmt.Errorf("fromBaseN: empty input")
	} // if

	base := uint64(len(model))

	for _, itor := range input {
		result *= base
		index := strings.IndexRune(model, itor)

		if index == -1 {
			return 0, fmt.Errorf("fromBaseN: invalid character %q in %q", itor, input)
		} // if

		result += uint64(index)
	} // for

	return result, nil
}
