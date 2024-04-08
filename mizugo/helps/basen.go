package helps

import (
	"fmt"
	"strings"
)

const (
	Base58 = "0123456789abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"                       // 58進制字串, 只使用0~9, a~z, A~Z 並且排除了oOlI這些容易混淆的字母
	Base80 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ._~+-*=?!|@#$%^&<>" // 80進制字串
)

// ToBase58 將uint64轉為58進制字串
func ToBase58(input uint64) string {
	return ToBaseN(Base58, input)
}

// FromBase58 將58進制字串轉為uint64
func FromBase58(input string) (uint64, error) {
	result, err := FromBaseN(Base58, input)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// ToBase80 將uint64轉為80進制字串
func ToBase80(input uint64) string {
	return ToBaseN(Base80, input)
}

// FromBase80 將80進制字串轉為uint64
func FromBase80(input string) (uint64, error) {
	result, err := FromBaseN(Base80, input)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// ToBaseN 將uint64轉為n進制字串
func ToBaseN(model string, input uint64) string {
	if input == 0 {
		return "0"
	} // if

	base := uint64(len(model))
	encoded := []string{}

	for input >= base {
		div, mod := input/base, input%base
		encoded = append([]string{string(model[mod])}, encoded...)
		input = div
	} // for

	if input > 0 {
		encoded = append([]string{string(model[input])}, encoded...)
	} // if

	return strings.Join(encoded, "")
}

// FromBaseN 將n進制字串轉為uint64
func FromBaseN(model, input string) (uint64, error) {
	base := uint64(len(model))
	value := uint64(0)

	for _, itor := range input {
		value *= base
		index := strings.IndexRune(model, itor)

		if index == -1 {
			return 0, fmt.Errorf("fromBaseN: invalid character %v in %v", itor, input)
		} // if

		value += uint64(index)
	} // for

	return value, nil
}
