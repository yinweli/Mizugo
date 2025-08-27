package testdata

import (
	"path/filepath"
	"runtime"
	"time"
)

const (
	Unknown         = "?????"                                     // 不明字串
	TestCount       = 100000                                      // 測試次數
	RedisTimeout    = time.Second                                 // redis逾時時間
	RedisIP         = "127.0.0.1:6379"                            // redis位址
	RedisURI        = "redisdb://127.0.0.1:6379/"                 // 有效redis連接字串
	RedisURIInvalid = "redisdb://127.0.0.1:10001/?dialTimeout=1s" // 無效redis連接字串
	MongoURI        = "mongodb://127.0.0.1:27017/"                // 有效mongo連接字串
	MongoURIInvalid = "mongodb://127.0.0.1:10001/?timeoutMS=1000" // 無效mongo連接字串
	TrialDirName    = "trial"                                     // 測試目錄名稱
	TrialFileName   = "trial.txt"                                 // 測試檔案名稱
)

// PathWork 取得測試工作路徑
func PathWork(work string) string {
	return filepath.Clean(filepath.Join(Root, work))
}

// PathEnv 取得測試環境路徑
func PathEnv(env string) string {
	return filepath.Clean(filepath.Join(Root, "env", env))
}

//nolint:gochecknoinits
func init() {
	_, file, _, _ := runtime.Caller(0) //nolint:dogsled
	Root = filepath.Clean(filepath.Dir(file))
	TrialDir = filepath.Join(Root, TrialDirName)
}

var Root string     // 測試根路徑
var TrialDir string // 測試目錄路徑
