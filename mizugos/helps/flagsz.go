package helps

import (
	"strings"
)

const (
	FlagszOnRune  = '1'                   // 開啟旗標字符
	FlagszOffRune = '0'                   // 關閉旗標字符
	flagszOn      = string(FlagszOnRune)  // 開啟旗標字串
	flagszOff     = string(FlagszOffRune) // 關閉旗標字串
)

// FlagszInit 初始化旗標字串
//
// size 指定旗標長度, flag 指定初始狀態
//   - flag=true  → 產生 size 個 '1'
//   - flag=false → 產生 size 個 '0'
//
// 若 size <= 0, 回傳空字串
func FlagszInit(size int, flag bool) string {
	if size <= 0 {
		return ""
	} // if

	if flag {
		return strings.Repeat(flagszOn, size)
	} else {
		return strings.Repeat(flagszOff, size)
	} // if
}

// FlagszSet 設定旗標字串指定索引位置的值
//
// 若 input 長度不足, 會自動補齊 '0' 直到該索引位置
//   - index < 0  → 直接回傳原字串
//   - flag=true  → 將該位置設為 '1'
//   - flag=false → 將該位置設為 '0'
func FlagszSet(input string, index int, flag bool) string {
	if index < 0 {
		return input
	} // if

	size := len(input)

	if size <= index {
		input += strings.Repeat(flagszOff, index-size+1)
	} // if

	return input[:index] + flagsz(flag) + input[index+1:]
}

// FlagszAdd 在旗標字串尾端新增一個旗標
func FlagszAdd(input string, flag bool) string {
	return input + flagsz(flag)
}

// FlagszAND 對兩個旗標字串做「逐位 AND 運算」
//
// 較短的字串會自動視為不足位補 '0'
func FlagszAND(input, other string) string {
	result := strings.Builder{}
	size := max(len(input), len(other))

	for i := 0; i < size; i++ {
		a := FlagszGet(input, i)
		b := FlagszGet(other, i)
		result.WriteString(flagsz(a && b))
	} // for

	return result.String()
}

// FlagszOR 對兩個旗標字串做「逐位 OR 運算」
//
// 較短的字串會自動視為不足位補 '0'
func FlagszOR(input, other string) string {
	result := strings.Builder{}
	size := max(len(input), len(other))

	for i := 0; i < size; i++ {
		a := FlagszGet(input, i)
		b := FlagszGet(other, i)
		result.WriteString(flagsz(a || b))
	} // for

	return result.String()
}

// FlagszXOR 對兩個旗標字串做「逐位 XOR 運算」
//
// 較短的字串會自動視為不足位補 '0'
func FlagszXOR(input, other string) string {
	result := strings.Builder{}
	size := max(len(input), len(other))

	for i := 0; i < size; i++ {
		a := FlagszGet(input, i)
		b := FlagszGet(other, i)
		result.WriteString(flagsz(a != b))
	} // for

	return result.String()
}

// FlagszGet 取得旗標字串在指定索引位置的狀態
//   - 若 index 在範圍內且為 '1' → 回傳 true
//   - 其他情況 → 回傳 false
func FlagszGet(input string, index int) bool {
	return index >= 0 && index < len(input) && input[index] == FlagszOnRune
}

// FlagszAny 判斷字串中是否至少有一個旗標為開啟 (至少一個 '1')
func FlagszAny(input string) bool {
	return strings.Contains(input, flagszOn)
}

// FlagszAll 判斷字串中是否所有旗標都為開啟 (全為 '1')
func FlagszAll(input string) bool {
	return strings.Contains(input, flagszOff) == false
}

// FlagszNone 判斷字串中是否所有旗標都為關閉 (全為 '0')
func FlagszNone(input string) bool {
	return strings.Contains(input, flagszOn) == false
}

// FlagszCount 計算字串中指定旗標的出現次數
func FlagszCount(input string, flag bool) int {
	return strings.Count(input, flagsz(flag))
}

// flagsz 取得旗標值代表的字串
func flagsz(flag bool) string {
	if flag {
		return flagszOn
	} else {
		return flagszOff
	} // if
}
