package main

import (
	"fmt"
	"runtime/debug"

	"github.com/yinweli/Mizugo/mizugos"
	"github.com/yinweli/Mizugo/mizugos/ctxs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
)

func main() {
	defer mizugos.Stop()
	defer func() {
		if cause := recover(); cause != nil && features.LogCrash != nil {
			features.LogCrash.Get().Error("crash").KV("stack", string(debug.Stack())).Error(fmt.Errorf("%s", cause)).EndFlush()
		} // if
	}()
	fmt.Println("test-client initialize")
	mizugos.Start()
	err := error(nil)

	if err = features.ConfigInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.LoggerInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.MetricsInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.PoolInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = entrys.AuthInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = entrys.JsonInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = entrys.ProtoInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = entrys.ProtoRavenInitialize(); err != nil {
		fmt.Println(err)
		return
	} // if

	fmt.Println("test-client running")
	ctx := ctxs.Get().WithCancel()

	for range ctx.Done() {
		// do nothing...
	} // for

	fmt.Println("test-client finalize")
}
