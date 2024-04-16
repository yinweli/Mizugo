package trials

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/otiai10/copy"
)

// Shell 測試框架處理, 依照以下流程執行
//   - 使用 Prepare 準備測試目錄
//   - 執行初始化處理
//   - 執行單元測試
//   - 執行結束處理
//   - 使用 Restore 還原測試目錄
//   - 回傳單元測試的結果編號, 可用此編號呼叫 os.Exit, 讓外部能夠獲得測試結果
func Shell(m *testing.M, initialize, finalize func(), work string, from ...string) (result int) {
	catalog := Prepare(work, from...)

	if initialize != nil {
		initialize()
		WaitTimeout() // 等待一下, 讓初始化處理有機會完成
	} // if

	if m != nil {
		result = m.Run()
	} // if

	if finalize != nil {
		finalize()
		WaitTimeout() // 等待一下, 讓結束處理有機會完成
	} // if

	Restore(catalog)
	return result
}

// Prepare 準備測試目錄, 依照以下流程執行
//   - 工作目錄改為 work 指定的路徑, 必須是絕對路徑
//   - 從 from 把資料複製過來
//   - 回傳用於還原的目錄資料
func Prepare(work string, from ...string) Catalog {
	if filepath.IsAbs(work) == false {
		panic("work must be absolute")
	} // if

	root, err := os.Getwd()

	if err != nil {
		panic(err)
	} // if

	if err = os.MkdirAll(work, os.ModePerm); err != nil {
		panic(err)
	} // if

	if err = os.Chdir(work); err != nil {
		panic(err)
	} // if

	for _, itor := range from {
		if err = copy.Copy(itor, "."); err != nil {
			panic(err)
		} // if
	} // for

	return Catalog{
		root: root,
		work: work,
	}
}

// Restore 還原測試目錄, 依照以下流程執行
//   - 工作目錄改為目錄資料中的原始路徑
//   - 刪除目錄資料中的工作路徑及其所有內容
func Restore(catalog Catalog) {
	if err := os.Chdir(catalog.root); err != nil {
		panic(err)
	} // if

	if err := os.RemoveAll(catalog.work); err != nil {
		panic(err)
	} // if
}

// Root 取得當前路徑
func Root() string {
	_, file, _, ok := runtime.Caller(1)

	if ok == false {
		panic("get root failed")
	} // if

	return filepath.Clean(filepath.Dir(file))
}

// Catalog 目錄資料
type Catalog struct {
	root string // 原始路徑
	work string // 工作路徑
}
