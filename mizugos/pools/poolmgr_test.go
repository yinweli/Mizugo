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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuitePoolmgr) SetupSuite() {
	this.Change("test-pools-poolmgr")
}

func (this *SuitePoolmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuitePoolmgr) TearDownTest() {
	this.GoLeak(this.T(), false) // 由於ants中有許多內部執行緒, 所以把這裡的執行緒洩漏檢查關閉
}

func (this *SuitePoolmgr) TestNewPoolmgr() {
	assert.NotNil(this.T(), NewPoolmgr())
}

func (this *SuitePoolmgr) TestInitialize() {
	target := NewPoolmgr()
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize(nil))
	assert.NotNil(this.T(), target.Initialize(nil)) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行

	target = NewPoolmgr()
	assert.Nil(this.T(), target.Initialize(&Config{
		Logger: &loggerTester{},
	}))
	target.Finalize()
}

func (this *SuitePoolmgr) TestSubmit() {
	target := NewPoolmgr()
	validNil := atomic.Bool{}
	validNilFunc := func() {
		validNil.Store(true)
	}
	assert.Nil(this.T(), target.Initialize(nil))
	target.Submit(validNilFunc)
	target.Finalize()
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validNil.Load())

	target = NewPoolmgr()
	validPool := atomic.Bool{}
	validPoolFunc := func() {
		validPool.Store(true)
	}
	assert.Nil(this.T(), target.Initialize(&Config{
		Logger: &loggerTester{},
	}))
	target.Submit(validPoolFunc)
	target.Finalize()
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), validPool.Load())

	target = NewPoolmgr()
	validFailed := atomic.Bool{}
	validFailedFunc := func() {
		validFailed.Store(true)
	}
	target.Submit(validFailedFunc)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), validFailed.Load())
}

func (this *SuitePoolmgr) TestStatus() {
	target := NewPoolmgr()
	assert.Equal(this.T(), Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	assert.Nil(this.T(), target.Initialize(nil))
	assert.Equal(this.T(), Stat{}, target.Status())
	target.Finalize()

	target = NewPoolmgr()
	assert.Nil(this.T(), target.Initialize(&Config{
		Logger: &loggerTester{},
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

type loggerTester struct {
}

func (this *loggerTester) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}
