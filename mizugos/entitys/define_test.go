package entitys

import (
	"sync/atomic"
)

// newModuleTester 建立模組測試器
func newModuleTester(moduleID ModuleID) *moduleTester {
	return &moduleTester{
		Module: NewModule(moduleID),
	}
}

// moduleTester 模組測試器
type moduleTester struct {
	*Module
	awake atomic.Int64
	start atomic.Int64
}

func (this *moduleTester) Awake() error {
	this.awake.Add(1)
	return nil
}

func (this *moduleTester) Start() error {
	this.start.Add(1)
	return nil
}
