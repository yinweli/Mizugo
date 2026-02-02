package trials

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/otiai10/copy"
)

// Shell 測試框架的統一入口, 用於包裝測試流程
//
// 執行順序
//   - 呼叫 Prepare 建立測試工作目錄
//   - 執行初始化處理(initialize)
//   - 執行測試
//   - 執行結束處理(finalize)
//   - 呼叫 Restore 還原測試目錄
//
// 回傳的結果編號, 可直接傳給 os.Exit 以回報測試結果
func Shell(m *testing.M, initialize, finalize func(), work string, from ...string) (result int) {
	catalog := Prepare(work, from...)
	defer func() {
		if r := recover(); r != nil {
			result = 1
		} // if

		Restore(catalog)
	}()

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

	return result
}

// Prepare 建立並切換至測試工作目錄
//
// 執行順序
//   - 驗證 work 必須是絕對路徑
//   - 建立目錄(若不存在)
//   - 切換至工作目錄
//   - 將 from 中的資料複製到工作目錄
//
// 回傳的 Catalog 將用於 Restore
func Prepare(work string, from ...string) Catalog {
	if filepath.IsAbs(work) == false {
		panic("shell: work must be absolute")
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

// Restore 還原目錄狀態
//
// 執行順序
//   - 切換回原始路徑
//   - 刪除工作目錄及其所有內容(若刪除失敗會重試最多10次)
func Restore(catalog Catalog) {
	const retry = 10
	const sleep = 100 * time.Millisecond

	if err := os.Chdir(catalog.root); err != nil {
		panic(err)
	} // if

	for i := 0; ; i++ {
		err := os.RemoveAll(catalog.work)

		if err == nil {
			return
		} // if

		if i > retry {
			panic(err)
		} // if

		time.Sleep(sleep)
	} // for
}

// Root 取得呼叫端檔案所在的目錄路徑; 若無法取得呼叫端資訊，會觸發 panic
func Root() string {
	_, file, _, ok := runtime.Caller(1)

	if ok == false {
		panic("shell: get root failed")
	} // if

	return filepath.Clean(filepath.Dir(file))
}

// Catalog 目錄資料
type Catalog struct {
	root string // 原始路徑
	work string // 工作路徑
}
