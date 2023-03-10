package defines

import (
	"fmt"
	"runtime/debug"

	"github.com/yinweli/Mizugo/mizugos"
)

// Crashlize 崩潰處理
func Crashlize(cause any) {
	mizugos.Error(LogCrash, "crash").KV("stack", string(debug.Stack())).EndError(fmt.Errorf("%s", cause))
}
