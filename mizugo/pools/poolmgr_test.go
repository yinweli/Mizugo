package pools

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestPoolmgr(t *testing.T) {
	suite.Run(t, new(SuitePoolmgr))
}

type SuitePoolmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuitePoolmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-pools-poolmgr")
}

func (this *SuitePoolmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuitePoolmgr) TearDownTest() {
	testdata.Leak(this.T(), false) // 由於ants中有許多內部執行緒, 所以把這裡的執行緒洩漏檢查關閉
}

func (this *SuitePoolmgr) TestPoolmgr() {
	target := NewPoolmgr()
	assert.NotNil(this.T(), target)
	config := &Config{
		Logger: func(format string, args ...any) {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		},
	}

	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize(config))
	assert.NotNil(this.T(), target.Initialize(config)) // 故意啟動兩次, 這次應該失敗
	assert.NotNil(this.T(), target.Initialize(nil))
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
}

func (this *SuitePoolmgr) TestSubmit() {
	target := NewPoolmgr()
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

	time.Sleep(testdata.Timeout)
	target.Submit(validFunc)
	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), int64(1), valid.Load())

	assert.Nil(this.T(), target.Initialize(config))
	time.Sleep(testdata.Timeout)
	target.Submit(validFunc)
	time.Sleep(testdata.Timeout)
	target.Finalize()
	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), int64(2), valid.Load())
}

func (this *SuitePoolmgr) TestStatus() {
	target := NewPoolmgr()
	assert.Equal(this.T(), Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	assert.NotNil(this.T(), target.Initialize(nil))
	assert.Equal(this.T(), Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	assert.Nil(this.T(), target.Initialize(&Config{
		Logger: func(format string, args ...any) {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		},
	}))
	assert.Equal(this.T(), Stat{Available: -1, Capacity: -1}, target.Status())
	target.Finalize()
}

func (this *SuitePoolmgr) TestConfig() {
	config := Config{}
	fmt.Println(config)
	assert.NotNil(this.T(), config.String())
}

func (this *SuitePoolmgr) TestStat() {
	stat := Stat{}
	fmt.Println(stat)
	assert.NotNil(this.T(), stat.String())
}
