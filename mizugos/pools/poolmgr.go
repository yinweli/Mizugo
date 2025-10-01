package pools

import (
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// NewPoolmgr 建立執行緒池管理器
func NewPoolmgr() *Poolmgr {
	return &Poolmgr{}
}

// Poolmgr 執行緒池管理器
type Poolmgr struct {
	logger ants.Logger  // 日誌物件
	pool   *ants.Pool   // 執行緒池
	lock   sync.RWMutex // 執行緒鎖
}

// Config 配置資料
type Config struct {
	Capacity     int                              `yaml:"capacity"`    // 執行緒池容量, 0表示容量無限
	Expire       time.Duration                    `yaml:"expire"`      // 執行緒逾時時間, 詳細說明請查看ants.Options.ExpiryDuration的說明
	PreAlloc     bool                             `yaml:"preAlloc"`    // 是否預先分配記憶體, 詳細說明請查看ants.Options.PreAlloc的說明
	Nonblocking  bool                             `yaml:"nonblocking"` // 是否在執行緒耗盡時阻塞Submit的執行, 詳細說明請查看ants.Options.Nonblocking的說明
	MaxBlocking  int                              `yaml:"maxBlocking"` // 最大阻塞執行緒數量, 0表示無限制, 詳細說明請查看ants.Options.MaxBlockingTasks的說明
	PanicHandler func(any)                        `yaml:"-" json:"-"`  // 失敗處理函式, 詳細說明請查看ants.Options.PanicHandler的說明
	Logger       func(format string, args ...any) `yaml:"-" json:"-"`  // 日誌函式
}

// Stat 狀態資料
type Stat struct {
	Running   int // 執行中的執行緒數量
	Available int // 未使用的執行緒數量
	Capacity  int // 執行緒數量上限
}

// Initialize 初始化處理
func (this *Poolmgr) Initialize(config *Config) (err error) {
	if config == nil {
		return fmt.Errorf("poolmgr initialize: config nil")
	} // if

	this.lock.Lock()
	defer this.lock.Unlock()

	if this.pool != nil {
		return fmt.Errorf("poolmgr initialize: already initialize")
	} // if

	this.logger = &Logger{logger: config.Logger}

	if this.pool, err = ants.NewPool(
		config.Capacity,
		ants.WithExpiryDuration(config.Expire),
		ants.WithPreAlloc(config.PreAlloc),
		ants.WithNonblocking(config.Nonblocking),
		ants.WithMaxBlockingTasks(config.MaxBlocking),
		ants.WithPanicHandler(config.PanicHandler),
		ants.WithLogger(this.logger),
	); err != nil {
		return fmt.Errorf("poolmgr initialize: %w", err)
	} // if

	return nil
}

// Finalize 結束處理
func (this *Poolmgr) Finalize() {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.pool != nil {
		this.pool.Release()
		this.pool = nil
		this.logger = nil
	} // if
}

// Submit 啟動執行緒
func (this *Poolmgr) Submit(task func()) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.pool == nil {
		go task()
		return
	} // if

	if err := this.pool.Submit(task); err != nil {
		this.logger.Printf("poolmgr submit: %v", err)
	} // if
}

// Status 獲得狀態資料
func (this *Poolmgr) Status() (result Stat) {
	this.lock.RLock()
	defer this.lock.RUnlock()

	if this.pool != nil {
		result = Stat{
			Running:   this.pool.Running(),
			Available: this.pool.Free(),
			Capacity:  this.pool.Cap(),
		}
	} // if

	return result
}

// Logger 日誌資料
type Logger struct {
	logger func(format string, args ...any)
}

// Printf 輸出日誌
func (this *Logger) Printf(format string, args ...any) {
	if this.logger != nil {
		this.logger(format, args...)
	} // if
}

func init() { //nolint:gochecknoinits
	DefaultPool = NewPoolmgr()
}

var DefaultPool *Poolmgr // 預設執行緒池管理器
