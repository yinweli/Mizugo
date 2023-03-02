package redmos

import (
	"context"
	"fmt"

	"github.com/yinweli/Mizugo/mizugos/ctxs"
)

// newMixed 建立混合資料庫
func newMixed(major *Major, minor *Minor) *Mixed {
	return &Mixed{
		major: major,
		minor: minor,
	}
}

// Mixed 混合資料庫, 內部用主要與次要資料庫實現混合命令, 包含以下功能
//   - 取得執行物件: 取得資料庫執行器
type Mixed struct {
	major *Major // 主要資料庫物件
	minor *Minor // 次要資料庫物件
}

// Submit 取得執行物件
func (this *Mixed) Submit(ctx ctxs.Ctx, tableName string) *Submit {
	return &Submit{
		ctx:   ctx,
		major: this.major.Submit(),
		minor: this.minor.Submit(tableName),
	}
}

// Submit 混合資料庫執行器, 執行命令時需要遵循以下流程
//   - 定義包含了 Behave 的行為結構 / 定義繼承了 Behavior 的行為結構
//   - 建立行為資料, 並且填寫內容(例如索引值, 資料內容, 結果成員等)
//   - 取得執行物件
//   - 新增行為到執行物件中
//   - 執行命令 Exec
//   - 檢查執行命令結果以及進行後續處理
//
// 目前已經實作了幾個預設行為 Lock, Unlock, Get, Set, Index 可以幫助使用者作行為設計;
// 其中 Lock, Unlock 已經直接整合到 Submit 提供的函式
type Submit struct {
	ctx      ctxs.Ctx    // ctx物件
	major    MajorSubmit // 主要執行物件
	minor    MinorSubmit // 次要執行物件
	behavior []Behavior  // 行為列表
}

// Add 新增行為
func (this *Submit) Add(behavior Behavior) *Submit {
	behavior.Initialize(this.ctx, this.major, this.minor)
	this.behavior = append(this.behavior, behavior)
	return this
}

// Lock 新增鎖定行為
func (this *Submit) Lock(key string) *Submit {
	return this.Add(&Lock{Key: key, time: Timeout})
}

// Unlock 新增解鎖行為
func (this *Submit) Unlock(key string) *Submit {
	return this.Add(&Unlock{Key: key})
}

// Exec 執行命令, 執行命令時, 會以下列順序執行
//   - 執行所有行為的 Prepare 函式, 如果有錯誤就中止執行
//   - 執行主要資料庫的管線處理
//   - 執行所有行為的 Complete 函式, 如果有錯誤就中止執行
//   - 清除管線資料
func (this *Submit) Exec() error {
	for _, itor := range this.behavior {
		if err := itor.Prepare(); err != nil {
			return fmt.Errorf("submit exec prepare: %w", err)
		} // if
	} // for

	_, _ = this.major.Exec(this.ctx.Ctx())

	for _, itor := range this.behavior {
		if err := itor.Complete(); err != nil {
			return fmt.Errorf("submit exec complete: %w", err)
		} // if
	} // for

	this.behavior = nil
	return nil
}

// Behavior 行為介面, 當建立行為時, 需要實現此介面, 建議可以把 Behave 組合進行為結構中, 可以省去初始處理的實作;
// 設計行為時, 有以下的設計規範
//   - 主要資料庫的情況下: 由於主要資料庫會用管線機制執行, 因此通常會在 Prepare 新增管線命令, 然後在 Complete 檢查執行是否符合預期
//   - 次要資料庫的情況下: 由於次要資料庫會以常規方式執行, 因此通常 Prepare 不會有內容, 然後在 Complete 執行資料庫命令並且檢查執行是否符合預期
//   - 錯誤處理: 當資料庫失敗時才會回傳錯誤, 若是邏輯錯誤(例如資料不存在), 就不應該回傳錯誤, 而是把結果記錄下來提供外部使用
type Behavior interface {
	// Initialize 初始處理
	Initialize(ctx ctxs.Ctx, major MajorSubmit, minor MinorSubmit)

	// Prepare 準備處理
	Prepare() error

	// Complete 完成處理
	Complete() error
}

// Behave 行為資料
type Behave struct {
	ctx   ctxs.Ctx    // ctx物件
	major MajorSubmit // 主要執行物件
	minor MinorSubmit // 次要執行物件
}

// Initialize 初始處理
func (this *Behave) Initialize(ctx ctxs.Ctx, major MajorSubmit, minor MinorSubmit) {
	this.ctx = ctx
	this.major = major
	this.minor = minor
}

// Ctx 取得ctx物件
func (this *Behave) Ctx() context.Context {
	return this.ctx.Ctx()
}

// Major 取得主要執行物件
func (this *Behave) Major() MajorSubmit {
	return this.major
}

// Minor 取得次要執行物件
func (this *Behave) Minor() MinorSubmit {
	return this.minor
}
