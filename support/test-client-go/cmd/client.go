package main

import (
	"fmt"
	"runtime/debug"

	"github.com/yinweli/Mizugo/mizugo/ctxs"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/entrys"
	"github.com/yinweli/Mizugo/support/test-client-go/internal/features"
)

func main() {
	defer features.FinalizeLogger()
	defer func() {
		if cause := recover(); cause != nil && features.LogCrash != nil {
			features.LogCrash.Get().Error("crash").KV("stack", string(debug.Stack())).Error(fmt.Errorf("%s", cause)).EndFlush()
		} // if
	}()

	fmt.Println("test-client start")
	err := error(nil)

	if err = features.InitializeConfig(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.InitializeLogger(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.InitializeEntity(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.InitializeLabel(); err != nil {
		fmt.Println(err)
		return
	} // if

	if err = features.InitializeMetrics(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer features.FinalizeMetrics()

	if err = features.InitializeNet(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer features.FinalizeNet()

	if err = features.InitializePool(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer features.FinalizePool()

	if err = entrys.InitializeAuth(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer entrys.FinalizeAuth()

	if err = entrys.InitializeJson(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer entrys.FinalizeJson()

	if err = entrys.InitializeProto(); err != nil {
		fmt.Println(err)
		return
	} // if

	defer entrys.FinalizeProto()

	fmt.Println("test-client running")
	ctx := ctxs.Get().WithCancel()

	for range ctx.Done() {
		// do nothing...
	} // for

	fmt.Println("test-client shutdown")
}
