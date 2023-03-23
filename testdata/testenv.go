package testdata

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/otiai10/copy"
)

// Env 環境資料
type Env struct {
	original string // 原始路徑
	workpath string // 工作路徑
}

// EnvSetup 設定環境
func EnvSetup(env *Env, work string, data ...string) {
	original, err := os.Getwd()

	if err != nil {
		panic(err)
	} // if

	workpath := filepath.Join(rootpath, work)

	if err = os.MkdirAll(workpath, os.ModePerm); err != nil {
		panic(err)
	} // if

	for _, itor := range data {
		if err = copy.Copy(filepath.Join(envpath, itor), workpath); err != nil {
			panic(err)
		} // if
	} // for

	if err = os.Chdir(workpath); err != nil {
		panic(err)
	} // if

	env.original = original
	env.workpath = workpath
}

// EnvRestore 還原環境
func EnvRestore(env *Env) {
	if err := os.Chdir(env.original); err != nil {
		panic(err)
	} // if

	if err := os.RemoveAll(env.workpath); err != nil {
		panic(err)
	} // if
}

// PathRoot 取得根路徑
func PathRoot() string {
	return rootpath
}

// PathEnv 取得環境路徑
func PathEnv() string {
	return envpath
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
