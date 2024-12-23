package helps

// MapKey 取得映射索引列表
func MapKey[K comparable, V any](input map[K]V) (result []K) {
	for itor := range input {
		result = append(result, itor)
	} // for

	return result
}

// MapValue 取得映射資料列表
func MapValue[K comparable, V any](input map[K]V) (result []V) {
	for _, itor := range input {
		result = append(result, itor)
	} // for

	return result
}

// MapFlatten 取得映射展平列表
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
