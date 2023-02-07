package dbsyncs

import (
	"github.com/redis/go-redis/v9"
)

// dbsyncs
// + 主要介面設計, 以Redis為核心設計
// + 次要介面設計, 以關係式資料庫為核心設計
// + 反向程序, 從關係式資料庫取回資料
// + 索引與鎖定
// + 儲存與恢復(store & restore)

// dbsyncs interface
// + Lock / Unlock
// + Get Data   -> from MajorLayer
// + Edit Data ...
// + Store Data -> save MajorLayer & save MinorLayer

type Depositor[T any] interface {
	MajorLoad(key string) T

	MajorSave(key string, data T)

	MinorLoad(key string) T

	MinorSave(key string, data T)
}

// Read
//   User -> Redis -> User
// Write
//   User -> Redis -> User
//        -> Hards
// Obverse Access: 順序存取, redis -> database
// Reverse Access: 逆序存取, database -> redis

// redis類型與用途
// + string
//   單值儲存/取用, 結構(json)儲存/取用
//   obverse: database:table:key:value
//   reverse: key : value
// + hash(map)
//   結構儲存/取用
//   obverse: database:table:key list:value list
//   reverse: hash { keys : values }
// + set
//   不重複列表
// + sorted set
//   有序列表
// + pub/sub
//   事件通道

type account struct {
	accountID int
	name      string
	level     int
	exp       int
	gold      int
}

type bag1 struct {
	bag []item
}

type item struct {
	itemID int
	count  int
}

type bag2 struct {
	bag map[int]int // [itemID, count]
}

func testme() {
	_ = redis.NewUniversalClient(nil)
}
