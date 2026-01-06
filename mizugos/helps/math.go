package helps

import (
	"cmp"
)

// Min 回傳列表中的最小值
//
// 若 input 為空(nil 或長度為 0), 將直接回傳 fallback 以確保安全
// 此設計旨在取代標準庫 slices.Min, 解決其在空列表時會觸發 Panic 的問題
//
// 泛型 E 支援任何可排序的型別(cmp.Ordered), 包含自定義的底層型別(Underlying Types)
func Min[S ~[]E, E cmp.Ordered](input S, fallback E) (result E) {
	if len(input) < 1 {
		return fallback
	} // if

	result = input[0]

	for i := 1; i < len(input); i++ {
		result = min(result, input[i])
	} // for

	return result
}

// Max 回傳列表中的最大值
//
// 若 input 為空(nil 或長度為 0), 將直接回傳 fallback 以確保安全
// 此設計旨在取代標準庫 slices.Max, 解決其在空列表時會觸發 Panic 的問題
//
// 泛型 E 支援任何可排序的型別(cmp.Ordered), 包含自定義的底層型別(Underlying Types)
func Max[S ~[]E, E cmp.Ordered](input S, fallback E) (result E) {
	if len(input) < 1 {
		return fallback
	} // if

	result = input[0]

	for i := 1; i < len(input); i++ {
		result = max(result, input[i])
	} // for

	return result
}
