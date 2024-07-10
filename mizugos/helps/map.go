package helps

// MapToArray 取得映射的索引與資料列表
func MapToArray[K comparable, V any](input map[K]V) (rk []K, rv []V) {
	for k, v := range input {
		rk = append(rk, k)
		rv = append(rv, v)
	} // for

	return rk, rv
}
