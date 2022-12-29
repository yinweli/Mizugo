package metrics

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMetrics(t *testing.T) {
	suite.Run(t, new(SuiteMetrics))
}

type SuiteMetrics struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMetrics) SetupSuite() {
	this.Change("test-metrics-metrics")
}

func (this *SuiteMetrics) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMetrics) TearDownTest() {
	// 由於用於監控的http伺服器到最後也不會關閉, 所以只好把這裡的執行緒洩漏檢查關閉
	// goleak.VerifyNone(this.T())
}

func (this *SuiteMetrics) TestInitialize() {
	Initialize(8080, &Auth{
		Username: "username",
		Password: "password",
	})
	Finalize()
}

func (this *SuiteMetrics) TestNew() {
	assert.NotNil(this.T(), NewInt("int"))
	assert.NotNil(this.T(), NewFloat("float"))
	assert.NotNil(this.T(), NewString("string"))
	assert.NotNil(this.T(), NewMap("map"))
	assert.NotNil(this.T(), NewRuntime("runtime"))
}

func (this *SuiteMetrics) TestRuntime() {
	Initialize(8080, nil)
	target := NewRuntime("test")
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
	Finalize()
}

func BenchmarkAdd(b *testing.B) {
	Initialize(8080, nil)
	value, _ := rand.Int(rand.Reader, big.NewInt(1000))
	target := NewRuntime(strconv.FormatInt(value.Int64(), 10))

	for i := 0; i < b.N; i++ {
		target.Add(time.Second)
	} // for

	Finalize()
}

func BenchmarkString(b *testing.B) {
	Initialize(8080, nil)
	value, _ := rand.Int(rand.Reader, big.NewInt(1000))
	target := NewRuntime(strconv.FormatInt(value.Int64(), 10))
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)
	target.Add(time.Second)

	for i := 0; i < b.N; i++ {
		_ = target.String()
	} // for

	Finalize()
}
