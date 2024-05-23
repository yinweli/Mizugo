package helps

import (
	"os"
	"runtime/pprof"
)

// Pprof CPU剖析工具
type Pprof struct {
	file *os.File // 結果檔案
}

// Start 開始CPU剖析
func (this *Pprof) Start(path string) (err error) {
	if this.file, err = os.Create(path); err != nil {
		return Err(err)
	} // if

	if err = pprof.StartCPUProfile(this.file); err != nil {
		_ = this.file.Close()
		this.file = nil
		return Err(err)
	} // if

	return nil
}

// Stop 停止CPU剖析
func (this *Pprof) Stop() {
	if this.file != nil {
		pprof.StopCPUProfile()
		_ = this.file.Close()
	} // if
}
