package helps

import (
	"reflect"
)

// ReflectFieldValue 取得反射物件中指定欄位的值
func ReflectFieldValue[T any](value reflect.Value, name string) (result T, ok bool) {
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	} // if

	field := value.FieldByName(name)

	if field.IsValid() == false {
		return result, false
	} // if

	result, ok = field.Interface().(T)

	if ok == false {
		return result, false
	} // if

	return result, true
}
