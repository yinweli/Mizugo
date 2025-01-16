package redmos

import (
	"github.com/redis/go-redis/v9"
)

// QPop 彈出佇列行為, 以索引值與資料到主要/次要資料中取得佇列資料, 使用上有以下幾點須注意
//   - 需要事先建立好資料結構, 並填寫到泛型類型T中, 請不要填入指標類型
//   - 資料結構的成員都需要設定好`bson:xxxxx`屬性
//   - 執行前設定好 MajorEnable, MinorEnable
//   - 執行前設定好 Name 並且不能為空字串
//   - 執行前設定好 Data, 如果為nil, 則內部程序會自己建立
//   - 執行後可用 Data 來取得資料
type QPop[T any] struct {
	Behave                       // 行為物件
	MajorEnable bool             // 啟用主要資料庫
	MinorEnable bool             // 啟用次要資料庫
	Name        string           // 佇列名稱, 同時會是主要資料庫中的索引值與次要資料庫中的表格名稱
	Data        *T               // 資料物件
	cmd         *redis.StringCmd // 命令結果
}
