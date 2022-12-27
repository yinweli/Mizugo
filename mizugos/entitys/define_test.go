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
	awake   atomic.Int64
	start   atomic.Int64
	dispose atomic.Int64
	update  atomic.Int64
}

func (this *moduleTester) Awake() {
	this.awake.Add(1)
}

func (this *moduleTester) Start() {
	this.start.Add(1)
}

func (this *moduleTester) Dispose() {
	this.dispose.Add(1)
}

func (this *moduleTester) Update() {
	this.update.Add(1)
}
