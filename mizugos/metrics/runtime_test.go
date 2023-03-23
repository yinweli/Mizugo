package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestRuntime(t *testing.T) {
	suite.Run(t, new(SuiteRuntime))
}

type SuiteRuntime struct {
	suite.Suite
	testdata.Env
	port int
}

func (this *SuiteRuntime) SetupSuite() {
	testdata.EnvSetup(&this.Env, "test-metrics-runtime")
	this.port = 8080
}

func (this *SuiteRuntime) TearDownSuite() {
	testdata.EnvRestore(&this.Env)
}

func (this *SuiteRuntime) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteRuntime) TestRuntime() {
	metricsmgr := NewMetricsmgr()
	assert.Nil(this.T(), metricsmgr.Initialize(this.port))
	time.Sleep(testdata.Timeout) // 等待一下, 讓初始化有機會完成

	target := metricsmgr.NewRuntime("test")
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.String())
	target.Add(time.Second)
	target.Rec()()
	target.Add(time.Second)
	target.Rec()()
	target.Add(time.Second)
	target.Rec()()
	time.Sleep(time.Second * 2)
	assert.NotNil(this.T(), target.String())
	fmt.Println(target.String())

	metricsmgr.Finalize()
}

func BenchmarkRuntimeRec(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	_ = metricsmgr.Initialize(8080)

	target := metricsmgr.NewRuntime(utils.RandString(10))

	for i := 0; i < b.N; i++ {
		target.Rec()()
	} // for

	metricsmgr.Finalize()
}

func BenchmarkRuntimeString(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	_ = metricsmgr.Initialize(8080)

	target := metricsmgr.NewRuntime(utils.RandString(10))
	target.Rec()()
	target.Rec()()
	target.Rec()()
	target.Rec()()
	target.Rec()()

	for i := 0; i < b.N; i++ {
		_ = target.String()
	} // for

	metricsmgr.Finalize()
}
