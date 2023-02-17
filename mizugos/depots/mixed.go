package depots

import (
	"context"
	"fmt"
	"time"
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

// Runner 取得執行物件
func (this *Mixed) Runner(ctx context.Context, dbName, tableName string) *Runner {
	return &Runner{
		ctx:         ctx,
		majorRunner: this.major.Runner(),
		minorRunner: this.minor.Runner(dbName, tableName),
	}
}

// Runner 混合資料庫執行器, 執行命令時需要遵循以下流程
//   - 定義繼承了 Action 的行為結構
//   - 建立行為資料, 並且填寫內容(例如索引值, 資料內容等)
//   - 取得執行物件
//   - 新增行為到執行物件中
//   - 執行命令 Exec
//   - 檢查執行命令結果以及進行後續處理
type Runner struct {
	ctx         context.Context // context物件
	majorRunner MajorRunner     // 主要執行物件
	minorRunner MinorRunner     // 次要執行物件
	action      []Action        // 行為列表
}

// Add 新增行為
func (this *Runner) Add(action Action) *Runner {
	this.action = append(this.action, action)
	return this
}

// Lock 新增鎖定行為
func (this *Runner) Lock(key string) *Runner {
	return this.Add(&Lock{Key: key, time: time.Second * 30}) // 預設鎖定超時時間為30秒
}

// Unlock 新增解鎖行為
func (this *Runner) Unlock(key string) *Runner {
	return this.Add(&Unlock{Key: key})
}

// Exec 執行命令, 執行命令時, 會以下列順序執行
//   - 執行所有行為的 Prepare 函式, 如果有錯誤就中止執行
//   - 執行主要資料庫的管線處理
//   - 執行所有行為的 Result 函式, 如果有錯誤就中止執行
func (this *Runner) Exec() error {
	for _, itor := range this.action {
		if err := itor.Prepare(this.ctx, this.majorRunner, this.minorRunner); err != nil {
			return fmt.Errorf("runner exec: %w", err)
		} // if
	} // for

	if _, err := this.majorRunner.Exec(this.ctx); err != nil {
		return fmt.Errorf("runner exec: %w", err)
	} // if

	for _, itor := range this.action {
		if err := itor.Result(); err != nil {
			return fmt.Errorf("runner exec: %w", err)
		} // if
	} // for

	return nil
}

// Action 行為介面, 在進行行為設計時, 針對主要/次要資料庫, 有以下的設計規範
//   - 主要資料庫的情況下: 由於主要資料庫會用管線機制執行, 因此通常會在 Prepare 新增管線命令, 然後在 Result 檢查執行結果是否符合預期
//   - 次要資料庫的情況下: 由於次要資料庫會以常規方式執行, 因此通常 Prepare 不會有內容, 然後在 Result 執行資料庫命令並且檢查執行結果是否符合預期
type Action interface {
	// Prepare 前置處理
	Prepare(ctx context.Context, majorRunner MajorRunner, minorRunner MinorRunner) error

	// Result 結果處理
	Result() error
}
