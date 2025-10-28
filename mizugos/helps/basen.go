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
	return ToBaseN(input, Base58)
}

// FromBase58 將 58 進制字串轉為 uint64
func FromBase58(input string) (result uint64, err error) {
	result, err = FromBaseN(input, Base58, base58Rank)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// LessBase58 比較 58 進制字串
func LessBase58(a, b string) bool {
	return LessBaseN(a, b, Base58, base58Rank)
}

// ToBase80 將 uint64 轉為 80 進制字串
func ToBase80(input uint64) string {
	return ToBaseN(input, Base80)
}

// FromBase80 將 80 進制字串轉為 uint64
func FromBase80(input string) (result uint64, err error) {
	result, err = FromBaseN(input, Base80, base80Rank)

	if err != nil {
		return 0, err
	} // if

	return result, nil
}

// LessBase80 比較 80 進制字串
func LessBase80(a, b string) bool {
	return LessBaseN(a, b, Base80, base80Rank)
}

// ToBaseN 將 uint64 轉為自訂進制字串, model 為字元表, 字元數量決定進制的基數, 字元表中不能有重複字元
func ToBaseN(input uint64, model string) string {
	if input == 0 {
		return model[:1]
	} // if

	base := uint64(len(model))
	result := [64]byte{}
	i := len(result)

	for input > 0 {
		i--
		result[i] = model[int(input%base)]
		input /= base
	} // for

	return string(result[i:])
}

// FromBaseN 將自訂進制字串轉為 uint64, model 為字元表, rank 為排序表
func FromBaseN(input, model string, rank *[256]int) (result uint64, err error) {
	if input == "" {
		return 0, fmt.Errorf("FromBaseN: empty input")
	} // if

	base := uint64(len(model))
	limit := ^uint64(0)
	input = strings.TrimLeft(input, model[:1])

	for i := 0; i < len(input); i++ {
		c := input[i]
		r := rank[c]

		if r < 0 { // 無效字元
			return 0, fmt.Errorf("FromBaseN: invalid character %q in %q", c, input)
		} // if

		u := uint64(r)

		if result > (limit-u)/base { // 溢位保護
			return 0, fmt.Errorf("FromBaseN: overflow in %q", input)
		} // if

		result = result*base + u
	} // for

	return result, nil
}

// LessBaseN 自訂進制字串排序, model 為字元表, rank 為排序表
func LessBaseN(a, b, model string, rank *[256]int) bool {
	zero := model[:1]
	a = strings.TrimLeft(a, zero)
	b = strings.TrimLeft(b, zero)
	la, lb := len(a), len(b)

	if la != lb { // 位數多者數值大
		return la < lb
	} // if

	for i := 0; i < la; i++ {
		ra, rb := rank[a[i]], rank[b[i]]

		if ra < 0 || rb < 0 { // 非法字元時退回字典序
			return a < b
		} // if

		if ra != rb {
			return ra < rb
		} // if
	} // for

	return false
}

// RankBaseN 建立自訂進制排序表, model 為字元表
func RankBaseN(model string) *[256]int {
	if len(model) < 2 {
		panic(fmt.Sprintf("RankBaseN: invalid model %q", model))
	} // if

	result := [256]int{}

	for i := range result {
		result[i] = -1
	} // for

	for i := 0; i < len(model); i++ {
		c := model[i]

		if result[c] >= 0 {
			panic(fmt.Sprintf("RankBaseN: duplicate %q in model %q", c, model))
		} // if

		result[c] = i
	} // for

	return &result
}

func init() { //nolint:gochecknoinits
	base58Rank = RankBaseN(Base58)
	base80Rank = RankBaseN(Base80)
}

// base58Rank 58 進制排序表
var base58Rank *[256]int

// base80Rank 80 進制排序表
var base80Rank *[256]int
