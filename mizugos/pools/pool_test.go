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

func TestPool(t *testing.T) {
	suite.Run(t, new(SuitePool))
}

type SuitePool struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuitePool) SetupSuite() {
	this.Change("test-pools-pool")
}

func (this *SuitePool) TearDownSuite() {
	this.Restore()
}

func (this *SuitePool) TearDownTest() {
	// 由於ants中有許多內部執行緒, 所以把這裡的執行緒洩漏檢查關閉
	// goleak.VerifyNone(this.T())
}

func (this *SuitePool) TestInitialize() {
	assert.Nil(this.T(), Initialize(Config{Logger: &loggerTester{}}))
	assert.NotNil(this.T(), Initialize(Config{}))
	Finalize()
}

func (this *SuitePool) TestSubmit() {
	poolLess := atomic.Bool{}
	assert.Nil(this.T(), Submit(func() {
		poolLess.Store(true)
	}))

	assert.Nil(this.T(), Initialize(Config{}))

	poolUsed := atomic.Bool{}
	assert.Nil(this.T(), Submit(func() {
		poolUsed.Store(true)
	}))

	Finalize()

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), poolLess.Load())
	assert.True(this.T(), poolUsed.Load())
}

func (this *SuitePool) TestStatus() {
	assert.Equal(this.T(), Stat{}, Status())
	assert.Nil(this.T(), Initialize(Config{}))
	assert.Equal(this.T(), Stat{Available: -1, Capacity: -1}, Status())
	Finalize()
}

func (this *SuitePool) TestConfig() {
	config := Config{}
	fmt.Println(config)
	assert.NotNil(this.T(), config.String())
}

func (this *SuitePool) TestStat() {
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
