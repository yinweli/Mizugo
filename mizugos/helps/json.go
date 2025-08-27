package helps

import (
	"encoding/json"
)

// JsonString 將 JSON 轉為字串
func JsonString(data any) string {
	result, _ := json.Marshal(data)
	return string(result)
}
