package pools

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// TODO: poolmgr
// TODO: 把執行緒換成由ants執行
// TODO: 由於ants會有錯誤, 所以函式回傳值可能得要改改
// TODO: 統一看看管理器的初始化/結束處理, 以及初始化時, 設定資料的統一, yaml的標記

// NewPoolmgr 建立執行緒池管理器
func NewPoolmgr() *Poolmgr {
	return &Poolmgr{}
}

// Poolmgr 執行緒池管理器
type Poolmgr struct {
	pool   *ants.Pool   // 執行緒池
	config *Config      // 設置資料
	lock   sync.RWMutex // 執行緒鎖
}

// Initialize 初始化處理
func (this *Poolmgr) Initialize(config *Config) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.pool != nil {
		return fmt.Errorf("poolmgr initialize: already initialize")
	} // if

	ants.Release() // 關閉預設的執行緒池
	pool, err := ants.NewPool(config.Capacity,
		ants.WithExpiryDuration(config.ExpireDuration),
		ants.WithPreAlloc(config.PreAlloc),
		ants.WithNonblocking(config.Nonblocking),
		ants.WithMaxBlockingTasks(config.MaxBlocking),
		ants.WithPanicHandler(config.PanicHandler),
		ants.WithLogger(config.Logger),
	)

	if err != nil {
		return fmt.Errorf("poolmgr initialize: %w", err)
	} // if

	this.pool = pool
	this.config = config
	this.logf("poolmgr start: %v", config)
	return nil
}

// Finalize 結束處理
func (this *Poolmgr) Finalize() {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.pool == nil {
		return
	} // uf

	this.logf("poolmgr stop")
	_ = this.pool.ReleaseTimeout(this.config.ReleaseDuration)
	this.pool = nil
	this.config = nil
}

// Submit 啟動執行緒
func (this *Poolmgr) Submit(task func()) error {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.pool == nil {
		go task()
		return nil
	} // if

	if err := this.pool.Submit(task); err != nil {
		return fmt.Errorf("poolmgr submit: %w", err)
	} // if

	return nil
}

// Status 獲得狀態資料
func (this *Poolmgr) Status() Stat {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.pool == nil {
		return Stat{}
	} // if

	return Stat{
		Running:   pool.Running(),
		Available: pool.Free(),
		Capacity:  pool.Cap(),
	}
}

// logf 記錄日誌
func (this *Poolmgr) logf(format string, args ...interface{}) {
	if this.config.Logger != nil {
		this.config.Logger.Printf(format, args...)
	} // if
}

// Config 設置資料
type Config struct {
	Capacity        int               `yaml:"capacity"`        // 執行緒池容量, 0表示容量無限
	ExpireDuration  time.Duration     `yaml:"expireDuration"`  // 執行緒逾時時間, 詳細說明請查看ants.Options.ExpiryDuration的說明
	ReleaseDuration time.Duration     `yaml:"releaseDuration"` // 關閉逾時時間, 當執行緒池結束時會等待此時間後開始關閉
	PreAlloc        bool              `yaml:"preAlloc"`        // 是否預先分配記憶體, 詳細說明請查看ants.Options.PreAlloc的說明
	Nonblocking     bool              `yaml:"nonblocking"`     // 是否在執行緒耗盡時阻塞Submit的執行, 詳細說明請查看ants.Options.Nonblocking的說明
	MaxBlocking     int               `yaml:"maxBlocking"`     // 最大阻塞執行緒數量, 0表示無限制, 詳細說明請查看ants.Options.MaxBlockingTasks的說明
	PanicHandler    func(interface{}) `yaml:"-"`               // 失敗處理函式, 詳細說明請查看ants.Options.PanicHandler的說明
	Logger          ants.Logger       `yaml:"-"`               // 日誌物件, 詳細說明請查看ants.Options.Logger的說明
}

// String 取得字串
func (this Config) String() string {
	return utils.ExpvarStr([]utils.ExpvarStat{
		{Name: "capacity", Data: this.Capacity},
		{Name: "expireDuration", Data: this.ExpireDuration},
		{Name: "releaseDuration", Data: this.ReleaseDuration},
		{Name: "preAlloc", Data: this.PreAlloc},
		{Name: "nonblocking", Data: this.Nonblocking},
		{Name: "maxBlocking", Data: this.MaxBlocking},
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
