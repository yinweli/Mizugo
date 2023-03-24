package testdata

import (
	"time"
)

const (
	Unknown          = "?????"                                     // 不明字串
	RandStringLength = 10                                          // 隨機字串長度
	Timeout          = time.Millisecond * 200                      // 超時時間
	RedisTimeout     = time.Second                                 // redis超時時間
	RedisURI         = "redisdb://127.0.0.1:6379/"                 // 有效redis連接字串
	RedisURIInvalid  = "redisdb://127.0.0.1:10001/?dialTimeout=1s" // 無效redis連接字串
	MongoURI         = "mongodb://127.0.0.1:27017/"                // 有效mongo連接字串
	MongoURIInvalid  = "mongodb://127.0.0.1:10001/?timeoutMS=1000" // 無效mongo連接字串
)
