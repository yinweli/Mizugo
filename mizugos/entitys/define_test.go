package entitys

import (
	"fmt"
	"sync/atomic"
)

// newTestModule 建立模組測試器
func newTestModule(awake, start bool, moduleID ModuleID) *testModule {
	return &testModule{
		Module: NewModule(moduleID),
		awake:  awake,
		start:  start,
	}
}

// testModule 模組測試器
type testModule struct {
	*Module
	awake      bool
	start      bool
	awakeCount atomic.Int64
	startCount atomic.Int64
}

func (this *testModule) Awake() error {
	this.awakeCount.Add(1)

	if this.awake {
		return nil
	} else {
		return fmt.Errorf("awake failed")
	} // if
}

func (this *testModule) Start() error {
	this.startCount.Add(1)

	if this.start {
		return nil
	} else {
		return fmt.Errorf("start failed")
	} // if
}
