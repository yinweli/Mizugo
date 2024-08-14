package helps

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// NewStdColor 建立色彩輸出顯示工具
func NewStdColor(stdout, stderr io.Writer) *StdColor {
	return &StdColor{
		stdout: stdout,
		stderr: stderr,
		failed: false,
	}
}

// StdColor 色彩輸出顯示工具
type StdColor struct {
	stdout io.Writer // 標準輸出
	stderr io.Writer // 錯誤輸出
	failed bool      // 失敗旗標
}

// Out 輸出標準訊息, 會自動附加換行符號到字串尾端
func (this *StdColor) Out(format string, a ...any) *StdColor {
	color.New(color.FgGreen).FprintlnFunc()(this.stdout, fmt.Sprintf(format, a...))
	return this
}

// Outf 輸出標準訊息
func (this *StdColor) Outf(format string, a ...any) *StdColor {
	color.New(color.FgGreen).FprintfFunc()(this.stdout, format, a...)
	return this
}

// Outln 輸出標準訊息
func (this *StdColor) Outln(a ...any) *StdColor {
	color.New(color.FgGreen).FprintlnFunc()(this.stdout, a...)
	return this
}

// Err 輸出錯誤訊息, 會自動附加換行符號到字串尾端
func (this *StdColor) Err(format string, a ...any) *StdColor {
	color.New(color.FgRed).FprintlnFunc()(this.stderr, fmt.Sprintf(format, a...))
	this.failed = true
	return this
}

// Errf 輸出錯誤訊息
func (this *StdColor) Errf(format string, a ...any) *StdColor {
	color.New(color.FgRed).FprintfFunc()(this.stderr, format, a...)
	this.failed = true
	return this
}

// Errln 輸出錯誤訊息
func (this *StdColor) Errln(a ...any) *StdColor {
	color.New(color.FgRed).FprintlnFunc()(this.stderr, a...)
	this.failed = true
	return this
}

// Failed 取得失敗旗標, 當有輸出過錯誤訊息, 則失敗旗標為true, 並且不會變回false
func (this *StdColor) Failed() bool {
	return this.failed
}

// GetStdout 取得標準輸出物件
func (this *StdColor) GetStdout() io.Writer {
	return stdColorWriter(func(p []byte) (n int, err error) {
		this.Outf("%v", string(p))
		return len(p), nil
	})
}

// GetStderr 取得錯誤輸出物件
func (this *StdColor) GetStderr() io.Writer {
	return stdColorWriter(func(p []byte) (n int, err error) {
		this.Errf("%v", string(p))
		return len(p), nil
	})
}

// stdColorWriter 色彩寫入函式類型
type stdColorWriter func(p []byte) (n int, err error)

// Write 寫入訊息
func (this stdColorWriter) Write(p []byte) (n int, err error) {
	return this(p)
}
