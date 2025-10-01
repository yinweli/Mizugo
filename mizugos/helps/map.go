package helps

// MapKey 回傳輸入 map 的所有鍵 (key) 的列表
func MapKey[K comparable, V any](input map[K]V) (result []K) {
	for itor := range input {
		result = append(result, itor)
	} // for

	return result
}

// MapValue 回傳輸入 map 的所有值 (value) 的列表
func MapValue[K comparable, V any](input map[K]V) (result []V) {
	for _, itor := range input {
		result = append(result, itor)
	} // for

	return result
}

// MapFlatten 將輸入 map 展平成一個列表, 其中每個元素為 MapFlattenData (含 key 與 value)
func MapFlatten[K comparable, V any](input map[K]V) (result []MapFlattenData[K, V]) {
	for k, v := range input {
		result = append(result, MapFlattenData[K, V]{
			K: k,
			V: v,
		})
	} // for

	return result
}

// MapFlattenData 展平資料
type MapFlattenData[K comparable, V any] struct {
	K K // 索引
	V V // 資料
}
