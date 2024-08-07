package metrics

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/helps"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestRuntime(t *testing.T) {
	suite.Run(t, new(SuiteRuntime))
}

type SuiteRuntime struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteRuntime) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-metrics-runtime"))
}

func (this *SuiteRuntime) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteRuntime) TestRuntime() {
	port := 9101
	metricsmgr := NewMetricsmgr()
	assert.Nil(this.T(), metricsmgr.Initialize(port))
	trials.WaitTimeout() // 等待一下, 讓初始化有機會完成

	target := metricsmgr.NewRuntime("test")
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.String())
	target.Add(time.Second)
	target.Rec()()
	target.Add(time.Second)
	target.Rec()()
	target.Add(time.Second)
	target.Rec()()
	trials.WaitTimeout(time.Second * 2)
	assert.NotNil(this.T(), target.String())
	fmt.Println(target.String())

	metricsmgr.Finalize()
}

func BenchmarkRuntimeRec(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	_ = metricsmgr.Initialize(8080)

	target := metricsmgr.NewRuntime(helps.RandStringDefault())

	for i := 0; i < b.N; i++ {
		target.Rec()()
	} // for

	metricsmgr.Finalize()
}

func BenchmarkRuntimeString(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	_ = metricsmgr.Initialize(8080)

	target := metricsmgr.NewRuntime(helps.RandStringDefault())
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
