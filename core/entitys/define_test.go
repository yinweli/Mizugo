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
	awake   atomic.Bool
	start   atomic.Bool
	dispose atomic.Bool
	update  atomic.Bool
}

func (this *moduleTester) Awake() {
	this.awake.Store(true)
}

func (this *moduleTester) Start() {
	this.start.Store(true)
}

func (this *moduleTester) Dispose() {
	this.dispose.Store(true)
}

func (this *moduleTester) Update() {
	this.update.Store(true)
}

// newReactTester 建立反應測試器
func newReactTester() *reactTester {
	return &reactTester{}
}

// reactTester 反應測試器
type reactTester struct {
}

func (this *reactTester) Encode(_ any) (packet []byte, err error) {
	return nil, nil
}

func (this *reactTester) Decode(_ []byte) (message any, err error) {
	return nil, nil
}

func (this *reactTester) Receive(_ any) error {
	return nil
}

func (this *reactTester) Error(_ error) {
}
