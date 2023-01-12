package pools

import (
	"fmt"
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
	config := &Config{
		ReleaseDuration: time.Second,
		Logger:          &loggerTester{},
	}

	assert.Nil(this.T(), target.Initialize(config))
	assert.NotNil(this.T(), target.Initialize(config))
	target.Finalize()
}

type loggerTester struct {
}

func (this *loggerTester) Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	fmt.Printf("\n")
}
