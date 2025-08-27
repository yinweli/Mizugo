package helps

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// NewStdColor 建立色彩輸出顯示工具
func NewStdColor(stdout, stderr io.Writer) *StdColor {
	if stdout == nil {
		stdout = io.Discard
	} // if

	if stderr == nil {
		stderr = io.Discard
	} // if

	g := color.New(color.FgGreen)
	r := color.New(color.FgRed)
	return &StdColor{
		stdout: stdout,
		stderr: stderr,
		failed: false,
		outf:   g.FprintfFunc(),
		outln:  g.FprintlnFunc(),
		errf:   r.FprintfFunc(),
		errln:  r.FprintlnFunc(),
	}
}

// StdColor 色彩輸出顯示工具
//
// 用於統一輸出附帶色彩的標準訊息與錯誤訊息
//
// 標準輸出 (stdout) 預設為綠色, 錯誤輸出 (stderr) 預設為紅色
//
// 若傳入的 stdout/stderr 為 nil, 會自動替換為 io.Discard, 避免後續輸出時 panic
type StdColor struct {
	stdout io.Writer // 標準輸出
	stderr io.Writer // 錯誤輸出
	failed bool      // 失敗旗標
	outf   func(w io.Writer, format string, a ...interface{})
	outln  func(w io.Writer, a ...interface{})
	errf   func(w io.Writer, format string, a ...interface{})
	errln  func(w io.Writer, a ...interface{})
}

// Out 輸出標準訊息, 會自動附加換行符號到字串尾端
func (this *StdColor) Out(format string, a ...any) *StdColor {
	this.outln(this.stdout, fmt.Sprintf(format, a...))
	return this
}

// Outf 輸出標準訊息
func (this *StdColor) Outf(format string, a ...any) *StdColor {
	this.outf(this.stdout, format, a...)
	return this
}

// Outln 輸出標準訊息
func (this *StdColor) Outln(a ...any) *StdColor {
	this.outln(this.stdout, a...)
	return this
}

// Err 輸出錯誤訊息, 會自動附加換行符號到字串尾端
func (this *StdColor) Err(format string, a ...any) *StdColor {
	this.errln(this.stderr, fmt.Sprintf(format, a...))
	this.failed = true
	return this
}

// Errf 輸出錯誤訊息
func (this *StdColor) Errf(format string, a ...any) *StdColor {
	this.errf(this.stderr, format, a...)
	this.failed = true
	return this
}

// Errln 輸出錯誤訊息
func (this *StdColor) Errln(a ...any) *StdColor {
	this.errln(this.stderr, a...)
	this.failed = true
	return this
}

// Failed 取得失敗旗標, 當有輸出過錯誤訊息, 則失敗旗標為true, 並且不會變回false
func (this *StdColor) Failed() bool {
	return this.failed
}

// GetStdout 取得標準輸出物件
func (this *StdColor) GetStdout() io.Writer {
	return stdColorWriter(func(b []byte) (n int, err error) {
		this.Outf("%s", b)
		return len(b), nil
	})
}

// GetStderr 取得錯誤輸出物件
func (this *StdColor) GetStderr() io.Writer {
	return stdColorWriter(func(b []byte) (n int, err error) {
		this.Errf("%s", b)
		return len(b), nil
	})
}

// stdColorWriter 色彩寫入函式類型
type stdColorWriter func(b []byte) (n int, err error)

// Write 寫入訊息
func (this stdColorWriter) Write(b []byte) (n int, err error) {
	return this(b)
}
