package helps

import (
	"encoding/json"
)

// JsonString 將json轉為字串
func JsonString(data any) string {
	result, _ := json.Marshal(data)
	return string(result)
}
