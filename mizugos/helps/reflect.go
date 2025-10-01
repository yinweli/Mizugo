package helps

import (
	"reflect"
)

// ReflectFieldValue 取得反射物件中指定欄位的值
//   - 傳入一個 reflect.Value (通常是 struct 或 struct 指標), 以及欄位名稱
//   - 嘗試取出該欄位的值並轉換為泛型 T
//   - 成功則回傳 (值, true), 失敗則回傳 (零值, false)
func ReflectFieldValue[T any](value reflect.Value, name string) (result T, ok bool) {
	if value.IsValid() == false {
		return result, false
	} // if

	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return result, false
		} // if

		value = value.Elem()
	} // if

	if value.Kind() != reflect.Struct {
		return result, false
	} // if

	field := value.FieldByName(name)

	if field.IsValid() == false {
		return result, false
	} // if

	if field.CanInterface() == false {
		return result, false
	} // if

	result, ok = field.Interface().(T)
	return result, ok
}
