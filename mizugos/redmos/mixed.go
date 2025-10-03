package redmos

import (
	"context"
	"fmt"
)

// newMixed 建立混合資料庫
func newMixed(major *Major, minor *Minor) *Mixed {
	return &Mixed{
		major: major,
		minor: minor,
	}
}

// Mixed 混合資料庫
//
// 以主要資料庫(Major)與次要資料庫(Minor)為基礎的混合封裝, 提供三種使用模式:
//   - Submit: 以「行為 (Behavior)」序列化的雙段式提交流程
//   - Major: 取得 Major, 適合需要 Redis 操作的場景
//   - Minor: 取得 Minor, 適合需要 Mongo 操作的場景
//
// Mixed 不是執行緒安全, 若跨 goroutine 共用, 請由上層管理器(Redmomgr)負責同步保護
type Mixed struct {
	major *Major // 主要資料庫物件
	minor *Minor // 次要資料庫物件
}

// Submit 取得混合執行器
func (this *Mixed) Submit(ctx context.Context) *Submit {
	return &Submit{
		context: ctx,
		major:   this.major.Submit(),
		minor:   this.minor.Submit(),
	}
}

// Major 取得主要資料庫(Major)
func (this *Mixed) Major() *Major {
	return this.major
}

// Minor 取得次要資料庫(Minor)
func (this *Mixed) Minor() *Minor {
	return this.minor
}

// Submit 混合執行器
//
// 提供以「行為(Behavior)」序列化的雙段式提交流程:
//
//	Prepare(排入 Redis 命令) → 主要資料庫管線送出 → Complete(依 Redis 結果決定 Mongo 寫入) → 次要資料庫批次送出
//
// 使用流程:
//   - 定義並實作行為(Behavior), 並將 Behave 組合入行為結構
//   - 建立行為資料並填入必要資訊(索引, 資料, 結果欄位等)
//   - 呼叫 Mixed.Submit 取得混合執行器
//   - 呼叫 Add 將行為加入執行佇列
//   - 呼叫 Exec 提交行為
//   - 依行為的結果欄位做後續處理(非錯誤情況下, 邏輯結果請從行為讀取)
//
// Submit 已內建以下的預設行為:
//   - Lock, Unlock: 鎖定, 解鎖
//   - LockIf, UnlockIf: 依旗標決定是否鎖定, 解鎖
type Submit struct {
	context  context.Context // ctx物件
	major    MajorSubmit     // 主要執行物件
	minor    *MinorSubmit    // 次要執行物件
	behavior []Behavior      // 行為列表
}

// Add 新增行為
//
// 會對每個行為呼叫 Behavior.Initialize 後再加入行為列表中
// 行為便可存取 Behave.Ctx / Behave.Major / Behave.Minor
func (this *Submit) Add(behavior ...Behavior) *Submit {
	for _, itor := range behavior {
		itor.Initialize(this.context, this.major, this.minor)
	} // for

	this.behavior = append(this.behavior, behavior...)
	return this
}

// Lock 新增鎖定行為
func (this *Submit) Lock(key, token string) *Submit {
	return this.Add(&Lock{Key: key, Token: token, ttl: Timeout})
}

// Unlock 新增解鎖行為
func (this *Submit) Unlock(key, token string) *Submit {
	return this.Add(&Unlock{Key: key, Token: token})
}

// LockIf 視旗標決定是否新增鎖定行為
func (this *Submit) LockIf(key, token string, lock bool) *Submit {
	if lock {
		return this.Add(&Lock{Key: key, Token: token, ttl: Timeout})
	} // if

	return this
}

// UnlockIf 視旗標決定是否新增解鎖行為
func (this *Submit) UnlockIf(key, token string, unlock bool) *Submit {
	if unlock {
		return this.Add(&Unlock{Key: key, Token: token})
	} // if

	return this
}

// Exec 依固定順序進行雙段式提交操作
//
// 執行流程:
//   - 依序呼叫所有行為的 Behavior.Prepare
//   - 呼叫主要資料庫 MajorSubmit 的 Exec
//   - 依序呼叫所有行為的 Behavior.Complete
//   - 呼叫次要資料庫 MinorSubmit.Exec
//   - 清除行為列表
//
// 執行時任何一步錯誤即停止並回傳
func (this *Submit) Exec() error {
	for _, itor := range this.behavior {
		if err := itor.Prepare(); err != nil {
			return fmt.Errorf("submit queue prepare: %w", err)
		} // if
	} // for

	if len(this.behavior) > 0 {
		_, _ = this.major.Exec(this.context)
	} // if

	for _, itor := range this.behavior {
		if err := itor.Complete(); err != nil {
			return fmt.Errorf("submit queue complete: %w", err)
		} // if
	} // for

	if len(this.behavior) > 0 {
		if err := this.minor.Exec(this.context); err != nil {
			return fmt.Errorf("submit queue minor: %w", err)
		} // if
	} // if

	this.behavior = nil
	return nil
}

// Behavior 行為介面
//
// 設計規範:
//   - 主要資料庫(Redis): 通常在 Prepare 階段加入批次命令, 於 Complete 階段檢查結果是否符合預期
//   - 次要資料庫(Mongo): 通常在 Complete 階段執行批次命令, 所以 Prepare 階段用於檢查參數或是為空
//   - 錯誤處理: 只有「資料庫失敗」才回傳錯誤; 邏輯條件(如資料不存在)不應回傳錯誤, 而應記錄在行為的結果欄位, 讓呼叫端判讀
type Behavior interface {
	// Initialize 初始處理
	Initialize(ctx context.Context, major MajorSubmit, minor *MinorSubmit)

	// Prepare 準備處理
	Prepare() error

	// Complete 完成處理
	Complete() error
}

// Behave 行為資料
//
// 提供行為的共用欄位與存取器, 省去 Initialize 程式碼, 須以組合方式內嵌到行為結構
type Behave struct {
	context context.Context // ctx物件
	major   MajorSubmit     // 主要執行物件
	minor   *MinorSubmit    // 次要執行物件
}

// Initialize 初始處理
func (this *Behave) Initialize(ctx context.Context, major MajorSubmit, minor *MinorSubmit) {
	this.context = ctx
	this.major = major
	this.minor = minor
}

// Ctx 取得 ctx 物件
func (this *Behave) Ctx() context.Context {
	return this.context
}

// Major 取得主要執行物件
func (this *Behave) Major() MajorSubmit {
	return this.major
}

// Minor 取得次要執行物件
func (this *Behave) Minor() *MinorSubmit {
	return this.minor
}
