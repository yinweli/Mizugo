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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteRuntime) SetupSuite() {
	this.Change("test-metrics-runtime")
}

func (this *SuiteRuntime) TearDownSuite() {
	this.Restore()
}

func (this *SuiteRuntime) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteRuntime) TestRuntime() {
	metricsmgr := NewMetricsmgr()
	metricsmgr.Initialize(8080)

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
	metricsmgr.Initialize(8080)

	target := metricsmgr.NewRuntime(utils.RandString(10))

	for i := 0; i < b.N; i++ {
		target.Rec()()
	} // for

	metricsmgr.Finalize()
}

func BenchmarkRuntimeString(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	metricsmgr.Initialize(8080)

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
