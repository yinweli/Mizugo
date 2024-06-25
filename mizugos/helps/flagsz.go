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
func FlagszInit(size int32, flag bool) string {
	if flag {
		return strings.Repeat(flagszOn, int(size))
	} else {
		return strings.Repeat(flagszOff, int(size))
	} // if
}

// FlagszSet 設定旗標, 若是旗標字串長度不足會自動填補
func FlagszSet(input string, index int32, flag bool) string {
	flagMax := int32(len(input))

	if flagMax <= index {
		input += strings.Repeat(flagszOff, int(index-flagMax+1))
	} // if

	return input[:index] + flagsz(flag) + input[index+1:]
}

// FlagszAdd 新增旗標, 新的旗標位於尾端
func FlagszAdd(input string, flag bool) string {
	return input + flagsz(flag)
}

// FlagszAND 對旗標字串做AND運算
func FlagszAND(input, other string) string {
	result := strings.Builder{}
	inputLen := len(input)
	otherLen := len(other)
	maxLen := inputLen

	if otherLen > inputLen {
		maxLen = otherLen
	} // if

	for i := 0; i < maxLen; i++ {
		a := FlagszGet(input, int32(i))
		b := FlagszGet(other, int32(i))
		result.WriteString(flagsz(a && b))
	} // for

	return result.String()
}

// FlagszOR 對旗標字串做OR運算
func FlagszOR(input, other string) string {
	result := strings.Builder{}
	inputLen := len(input)
	otherLen := len(other)
	maxLen := inputLen

	if otherLen > inputLen {
		maxLen = otherLen
	} // if

	for i := 0; i < maxLen; i++ {
		a := FlagszGet(input, int32(i))
		b := FlagszGet(other, int32(i))
		result.WriteString(flagsz(a || b))
	} // for

	return result.String()
}

// FlagszXOR 對旗標字串做XOR運算
func FlagszXOR(input, other string) string {
	result := strings.Builder{}
	inputLen := len(input)
	otherLen := len(other)
	maxLen := inputLen

	if otherLen > inputLen {
		maxLen = otherLen
	} // if

	for i := 0; i < maxLen; i++ {
		a := FlagszGet(input, int32(i))
		b := FlagszGet(other, int32(i))
		result.WriteString(flagsz(a != b))
	} // for

	return result.String()
}

// FlagszGet 取得旗標
func FlagszGet(input string, index int32) bool {
	return index < int32(len(input)) && input[index] == FlagszOnRune
}

// FlagszAny 是否為任一旗標開啟
func FlagszAny(input string) bool {
	return strings.Contains(input, flagszOn)
}

// FlagszAll 是否為全部旗標開啟
func FlagszAll(input string) bool {
	return strings.Contains(input, flagszOff) == false
}

// FlagszNone 是否為全部旗標關閉
func FlagszNone(input string) bool {
	return strings.Contains(input, flagszOn) == false
}

// FlagszCount 取得旗標的出現數量
func FlagszCount(input string, flag bool) int32 {
	return int32(strings.Count(input, flagsz(flag)))
}

// flagsz 取得旗標值代表的字串
func flagsz(flag bool) string {
	if flag {
		return flagszOn
	} else {
		return flagszOff
	} // if
}
