package utils

import (
	"io"

	"github.com/fatih/color"
)

// NewStdColor 建立色彩輸出顯示工具
func NewStdColor(stdout, stderr io.Writer) *StdColor {
	return &StdColor{
		stdout: stdout,
		stderr: stderr,
	}
}

// StdColor 色彩輸出顯示工具
type StdColor struct {
	stdout io.Writer // 標準輸出
	stderr io.Writer // 錯誤輸出
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

// Errf 輸出錯誤訊息
func (this *StdColor) Errf(format string, a ...any) *StdColor {
	color.New(color.FgRed).FprintfFunc()(this.stderr, format, a...)
	return this
}

// Errln 輸出錯誤訊息
func (this *StdColor) Errln(a ...any) *StdColor {
	color.New(color.FgRed).FprintlnFunc()(this.stderr, a...)
	return this
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
