package pools

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestPoolmgr(t *testing.T) {
	suite.Run(t, new(SuitePoolmgr))
}

type SuitePoolmgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuitePoolmgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-pools-poolmgr"))
}

func (this *SuitePoolmgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuitePoolmgr) TestPoolmgr() {
	config := &Config{
		Logger: func(format string, args ...any) {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		},
	}
	target := NewPoolmgr()
	this.NotNil(target)
	target.Finalize() // 初始化前執行, 這次應該不執行
	this.Nil(target.Initialize(config))
	this.NotNil(target.Initialize(config)) // 故意啟動兩次, 這次應該失敗
	this.NotNil(target.Initialize(nil))
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
}

func (this *SuitePoolmgr) TestSubmit() {
	config := &Config{
		Logger: func(format string, args ...any) {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		},
	}
	valid := atomic.Int64{}
	validFunc := func() {
		valid.Add(1)
	}
	target := NewPoolmgr()
	trials.WaitTimeout()
	target.Submit(validFunc)
	trials.WaitTimeout()
	this.Equal(int64(1), valid.Load())

	_ = target.Initialize(config)
	trials.WaitTimeout()
	target.Submit(validFunc)
	trials.WaitTimeout()
	target.Finalize()
	trials.WaitTimeout()
	this.Equal(int64(2), valid.Load())
}

func (this *SuitePoolmgr) TestStatus() {
	target := NewPoolmgr()
	this.Equal(Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	_ = target.Initialize(nil)
	this.Equal(Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	_ = target.Initialize(&Config{
		Logger: func(format string, args ...any) {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		},
	})
	this.Equal(Stat{Available: -1, Capacity: -1}, target.Status())
	target.Finalize()
}

func (this *SuitePoolmgr) TestLogger() {
	target := false
	(&Logger{logger: func(format string, args ...any) { target = true }}).Printf("")
	this.True(target)
}
