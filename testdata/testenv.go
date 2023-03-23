package testdata

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/otiai10/copy"
	"go.uber.org/goleak"
)

// TestEnv 測試環境
type TestEnv struct {
	// 通用測試

	Unknown string // 未知字串

	// 測試環境

	original string // 原始路徑
	workpath string // 工作路徑
}

// TBegin 開始測試
func (this *TestEnv) TBegin(work, data string) {
	// 初始化通用測試

	this.Unknown = "?????"

	// 初始化測試環境

	original, err := os.Getwd()

	if err != nil {
		panic(err)
	} // if

	workpath := filepath.Join(rootpath, work)

	if err = os.MkdirAll(workpath, os.ModePerm); err != nil {
		panic(err)
	} // if

	datapath := filepath.Join(envpath, data)

	if err = copy.Copy(datapath, workpath); err != nil {
		panic(err)
	} // if

	if err = os.Chdir(workpath); err != nil {
		panic(err)
	} // if

	this.original = original
	this.workpath = workpath
}

// TFinal 結束測試
func (this *TestEnv) TFinal() {
	if err := os.Chdir(this.original); err != nil {
		panic(err)
	} // if

	if err := os.RemoveAll(this.workpath); err != nil {
		panic(err)
	} // if
}

// TLeak 執行洩漏測試, 測試是否有執行緒未被關閉, 但是會有誤判的狀況, 預設關閉;
// 可以通過在 init 函式中把 leakTest 設為 true 來開啟
func (this *TestEnv) TLeak(t goleak.TestingT, run bool) {
	if leakTest && run {
		goleak.VerifyNone(t)
	} // if
}

// TCompareFile 比對檔案內容, 預期資料來自位元陣列
func (this *TestEnv) TCompareFile(path string, expected []byte) bool {
	if actual, err := os.ReadFile(path); err == nil {
		return string(expected) == string(actual)
	} // if

	return false
}

// RootPath 取得根路徑
func RootPath() string {
	return rootpath
}

// EnvPath 取得環境路徑
func EnvPath() string {
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

	leakTest = false // 這裡用來開啟/關閉洩漏測試旗標
}

var rootpath string // 根路徑
var envpath string  // 環境路徑
var leakTest bool   // 洩漏測試旗標
