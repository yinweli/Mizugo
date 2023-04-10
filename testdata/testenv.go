package testdata

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/otiai10/copy"
)

// EnvSetup 設定測試環境, 依照以下流程執行
//   - 工作目錄改為使用者指定的目錄
//   - 從使用者指定的目錄把測試資料複製過來
func EnvSetup(work string, data ...string) Env {
	original, err := os.Getwd()

	if err != nil {
		panic(err)
	} // if

	workpath := filepath.Join(rootpath, work)

	if err = os.MkdirAll(workpath, os.ModePerm); err != nil {
		panic(err)
	} // if

	if err = os.Chdir(workpath); err != nil {
		panic(err)
	} // if

	for _, itor := range data {
		if err = copy.Copy(filepath.Join(envpath, itor), "."); err != nil {
			panic(err)
		} // if
	} // for

	return Env{
		original: original,
		workpath: workpath,
	}
}

// EnvRestore 還原環境
func EnvRestore(env Env) {
	if err := os.Chdir(env.original); err != nil {
		panic(err)
	} // if

	if err := os.RemoveAll(env.workpath); err != nil {
		panic(err)
	} // if
}

// Env 環境資料
type Env struct {
	original string // 原始路徑
	workpath string // 工作路徑
}

func init() {
	_, file, _, ok := runtime.Caller(0)

	if ok == false {
		panic("get rootpath failed")
	} // if

	rootpath = filepath.Dir(file)
	envpath = filepath.Join(rootpath, "env")

	// 如果env資料夾不存在, 就建立一個, 免得測試時拋出錯誤
	if _, err := os.Stat(envpath); os.IsNotExist(err) {
		if err = os.MkdirAll(envpath, os.ModePerm); err != nil {
			panic(err)
		} // if
	} // if
}

var rootpath string // 根路徑
var envpath string  // 環境路徑
