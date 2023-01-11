package pools

import (
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// Initialize 初始化處理
func Initialize(config Config) error {
	if pool != nil {
		return fmt.Errorf("pool initialize: already initialize")
	} // if

	p, err := ants.NewPool(config.Capacity,
		ants.WithExpiryDuration(config.Expire),
		ants.WithPreAlloc(config.PreAlloc),
		ants.WithMaxBlockingTasks(config.MaxBlockingTasks),
		ants.WithNonblocking(config.Nonblocking),
		ants.WithPanicHandler(config.PanicHandler),
		ants.WithLogger(config.Logger),
	)

	if err != nil {
		return fmt.Errorf("pool initialize: %w", err)
	} // if

	pool = p

	if config.Logger != nil {
		config.Logger.Printf("pool start: %v", config)
	} // if

	return nil
}

// Finalize 結束處理
func Finalize() {
	if pool != nil {
		pool.Release()
		pool = nil
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

// Config 選項資料
type Config struct {
	Capacity         int               // 執行緒池容量, 0表示容量無限
	Expire           time.Duration     // 執行緒逾時時間, 詳細說明請查看ants.Options.ExpiryDuration的說明
	PreAlloc         bool              // 是否預先分配記憶體, 詳細說明請查看ants.Options.PreAlloc的說明
	MaxBlockingTasks int               // 最大阻塞執行緒數量, 0表示無限制, 詳細說明請查看ants.Options.MaxBlockingTasks的說明
	Nonblocking      bool              // 是否在執行緒耗盡時阻塞Submit的執行, 詳細說明請查看ants.Options.Nonblocking的說明
	PanicHandler     func(interface{}) // 失敗處理函式, 詳細說明請查看ants.Options.PanicHandler的說明
	Logger           ants.Logger       // 日誌物件, 詳細說明請查看ants.Options.Logger的說明
}

// String 取得字串
func (this Config) String() string {
	return utils.ExpvarStr([]utils.ExpvarStat{
		{Name: "capacity", Data: this.Capacity},
		{Name: "expire", Data: this.Expire},
		{Name: "preAlloc", Data: this.PreAlloc},
		{Name: "maxBlockingTasks", Data: this.MaxBlockingTasks},
		{Name: "nonblocking", Data: this.Nonblocking},
	})
}

// Stat 狀態資料
type Stat struct {
	Running   int // 執行中的執行緒數量
	Available int // 未使用的執行緒數量
	Capacity  int // 執行緒數量上限
}

// String 取得字串
func (this Stat) String() string {
	return utils.ExpvarStr([]utils.ExpvarStat{
		{Name: "running", Data: this.Running},
		{Name: "available", Data: this.Available},
		{Name: "capacity", Data: this.Capacity},
	})
}

var pool *ants.Pool // 執行緒池
