package miscs

import (
	"time"

	"github.com/yinweli/Mizugo/v2/mizugos"
)

// GenerateConnection 產生連線
func GenerateConnection(internal time.Duration, count, batch int, done func()) {
	mizugos.Pool.Submit(func() {
		timeout := time.NewTicker(internal)

		for range timeout.C {
			for i := 0; count > 0 && i < batch; i++ {
				done()
				count--
			} // for

			if count <= 0 {
				timeout.Stop()
				return
			} // if
		} // for
	})
}
