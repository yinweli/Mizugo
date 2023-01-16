package entitys

import (
	"fmt"
	"sync/atomic"
)

// newModuleTester 建立模組測試器
func newModuleTester(awake, start bool, moduleID ModuleID) *moduleTester {
	return &moduleTester{
		Module: NewModule(moduleID),
		awake:  awake,
		start:  start,
	}
}

// moduleTester 模組測試器
type moduleTester struct {
	*Module
	awake      bool
	start      bool
	awakeCount atomic.Int64
	startCount atomic.Int64
}

func (this *moduleTester) Awake() error {
	this.awakeCount.Add(1)

	if this.awake {
		return nil
	} else {
		return fmt.Errorf("awake failed")
	} // if
}

func (this *moduleTester) Start() error {
	this.startCount.Add(1)

	if this.start {
		return nil
	} else {
		return fmt.Errorf("start failed")
	} // if
}
