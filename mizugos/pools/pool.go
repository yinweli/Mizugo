package pools

import (
	"fmt"

	"github.com/panjf2000/ants/v2"
)

// Initialize 初始化處理
func Initialize(config Config) error {
	if pool != nil {
		return fmt.Errorf("pool initialize: already initialize")
	} // if

	ants.Release() // 關閉預設的執行緒池
	p, err := ants.NewPool(config.Capacity,
		ants.WithExpiryDuration(config.ExpireDuration),
		ants.WithPreAlloc(config.PreAlloc),
		ants.WithNonblocking(config.Nonblocking),
		ants.WithMaxBlockingTasks(config.MaxBlocking),
		ants.WithPanicHandler(config.PanicHandler),
		ants.WithLogger(config.Logger),
	)

	if err != nil {
		return fmt.Errorf("pool initialize: %w", err)
	} // if

	pool = p
	poolConfig = config
	logf("pool start: %v", config)
	return nil
}

// Finalize 結束處理
func Finalize() {
	logf("pool stop")

	if pool != nil {
		_ = pool.ReleaseTimeout(poolConfig.ReleaseDuration)
		pool = nil
		poolConfig = Config{}
	} // if
}

// Submit 啟動執行緒
func Submit(task func()) error {
	if pool != nil {
		if err := pool.Submit(task); err != nil {
			return fmt.Errorf("pool submit: %w", err)
		} // if
	} // if

	go task()
	return nil
}

// Status 獲得狀態資料
func Status() Stat {
	if pool != nil {
		return Stat{
			Running:   pool.Running(),
			Available: pool.Free(),
			Capacity:  pool.Cap(),
		}
	} // if

	return Stat{}
}

// logf 記錄日誌
func logf(format string, args ...interface{}) {
	if poolConfig.Logger != nil {
		poolConfig.Logger.Printf(format, args...)
	} // if
}

var pool *ants.Pool   // 執行緒池
var poolConfig Config // 執行緒池設置
