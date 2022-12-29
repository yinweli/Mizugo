package metrics

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestRuntime(t *testing.T) {
	suite.Run(t, new(SuiteRuntime))
}

type SuiteRuntime struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteRuntime) SetupSuite() {
	this.Change("test-metrics-runtime")
}

func (this *SuiteRuntime) TearDownSuite() {
	this.Restore()
}

func (this *SuiteRuntime) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteRuntime) TestRuntime() {
	metricsmgr := NewMetricsmgr()
	metricsmgr.Initialize(8080, nil)

	target := metricsmgr.NewRuntime("test")
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.String())
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Rec()()
	time.Sleep(time.Second * 2)
	assert.NotNil(this.T(), target.String())
	fmt.Println(target.String())

	metricsmgr.Finalize()
}

func BenchmarkRuntimeAdd(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	metricsmgr.Initialize(8080, nil)

	target := metricsmgr.NewRuntime(randString(10))

	for i := 0; i < b.N; i++ {
		target.Add(time.Second)
	} // for

	metricsmgr.Finalize()
}

func BenchmarkRuntimeString(b *testing.B) {
	metricsmgr := NewMetricsmgr()
	metricsmgr.Initialize(8080, nil)

	target := metricsmgr.NewRuntime(randString(10))
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)

	for i := 0; i < b.N; i++ {
		_ = target.String()
	} // for

	metricsmgr.Finalize()
}

func randString(count int) string {
	letter := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	generator, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letter))))
	builder := bytes.Buffer{}

	for i := 0; i < count; i++ {
		index := int(generator.Int64())
		builder.WriteByte(letter[index])
	} // for

	return builder.String()
}
