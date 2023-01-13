package pools

import (
	"fmt"
	"time"

	"github.com/panjf2000/ants/v2"

	"github.com/yinweli/Mizugo/mizugos/utils"
)

// TODO: 把執行緒換成由ants執行
// TODO: 由於ants會有錯誤, 所以函式回傳值可能得要改改

// NewPoolmgr 建立執行緒池管理器
func NewPoolmgr() *Poolmgr {
	return &Poolmgr{}
}

// Poolmgr 執行緒池管理器
type Poolmgr struct {
	pool *ants.Pool     // 執行緒池
	once utils.SyncOnce // 單次執行物件
}

// Initialize 初始化處理
func (this *Poolmgr) Initialize(config *Config) (err error) {
	if this.once.Done() {
		return fmt.Errorf("poolmgr initialize: already initialize")
	} // if

	this.once.Do(func() {
		if config == nil {
			return
		} // if

		this.pool, err = ants.NewPool(
			config.Capacity,
			ants.WithExpiryDuration(config.Expire),
			ants.WithPreAlloc(config.PreAlloc),
			ants.WithNonblocking(config.Nonblocking),
			ants.WithMaxBlockingTasks(config.MaxBlocking),
			ants.WithPanicHandler(config.PanicHandler),
			ants.WithLogger(config.Logger),
		)

		if err != nil {
			err = fmt.Errorf("poolmgr initialize: %w", err)
		} // if
	})

	return err
}

// Finalize 結束處理
func (this *Poolmgr) Finalize() {
	if this.once.Done() == false {
		return
	} // if

	if this.pool != nil {
		this.pool.Release()
	} // if
}

// Submit 啟動執行緒
func (this *Poolmgr) Submit(task func()) error {
	if this.once.Done() == false {
		return fmt.Errorf("poolmgr submit: not initialize")
	} // if

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
	if this.once.Done() == false {
		return Stat{}
	} // if

	if this.pool == nil {
		return Stat{}
	} // if

	return Stat{
		Running:   this.pool.Running(),
		Available: this.pool.Free(),
		Capacity:  this.pool.Cap(),
	}
}

// Config 設置資料
type Config struct {
	Capacity     int               `yaml:"capacity"`    // 執行緒池容量, 0表示容量無限
	Expire       time.Duration     `yaml:"expire"`      // 執行緒逾時時間, 詳細說明請查看ants.Options.ExpiryDuration的說明
	PreAlloc     bool              `yaml:"preAlloc"`    // 是否預先分配記憶體, 詳細說明請查看ants.Options.PreAlloc的說明
	Nonblocking  bool              `yaml:"nonblocking"` // 是否在執行緒耗盡時阻塞Submit的執行, 詳細說明請查看ants.Options.Nonblocking的說明
	MaxBlocking  int               `yaml:"maxBlocking"` // 最大阻塞執行緒數量, 0表示無限制, 詳細說明請查看ants.Options.MaxBlockingTasks的說明
	PanicHandler func(interface{}) `yaml:"-"`           // 失敗處理函式, 詳細說明請查看ants.Options.PanicHandler的說明
	Logger       ants.Logger       `yaml:"-"`           // 日誌物件, 詳細說明請查看ants.Options.Logger的說明
}

// String 取得字串
func (this Config) String() string {
	return utils.ExpvarStr([]utils.ExpvarStat{
		{Name: "capacity", Data: this.Capacity},
		{Name: "expire", Data: this.Expire},
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

func init() { //nolint
	DefaultPool = NewPoolmgr()
}

var DefaultPool *Poolmgr // 預設執行緒池管理器
