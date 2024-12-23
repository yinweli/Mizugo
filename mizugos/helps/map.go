package helps

// MapFlatten 映射展平, 將映射的內容展平為列表
func MapFlatten[K comparable, V any](input map[K]V) (result []MapFlattenData[K, V]) {
	for k, v := range input {
		result = append(result, MapFlattenData[K, V]{
			K: k,
			V: v,
		})
	} // for

	return result
}

// MapFlattenData 映射展平資料
type MapFlattenData[K comparable, V any] struct {
	K K // 索引
	V V // 資料
}
