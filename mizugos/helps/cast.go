package helps

import (
	"fmt"
)

// CastPointer 指標轉換
func CastPointer[T any](input any) (output *T, err error) {
	if input == nil {
		return nil, fmt.Errorf("cast pointer: input nil")
	} // if

	pointer, ok := input.(*T)

	if ok == false {
		return nil, fmt.Errorf("cast pointer: type failed")
	} // if

	return pointer, nil
}
