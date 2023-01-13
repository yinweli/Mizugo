package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMetricsmgr(t *testing.T) {
	suite.Run(t, new(SuiteMetricsmgr))
}

type SuiteMetricsmgr struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	port int
}

func (this *SuiteMetricsmgr) SetupSuite() {
	this.Change("test-metrics-metricsmgr")
	this.port = 8080
}

func (this *SuiteMetricsmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMetricsmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMetricsmgr) TestInitialize() {
	target := NewMetricsmgr()
	target.Finalize() // 初始化前執行, 這次應該不執行
	assert.Nil(this.T(), target.Initialize(this.port))
	assert.NotNil(this.T(), target.Initialize(this.port)) // 故意啟動兩次, 這次應該失敗
	target.Finalize()
	target.Finalize() // 故意結束兩次, 這次應該不執行
}

func (this *SuiteMetricsmgr) TestNew() {
	target := NewMetricsmgr()
	assert.Nil(this.T(), target.Initialize(this.port))

	assert.NotNil(this.T(), target.NewInt("int"))
	assert.NotNil(this.T(), target.NewFloat("float"))
	assert.NotNil(this.T(), target.NewString("string"))
	assert.NotNil(this.T(), target.NewMap("map"))
	assert.NotNil(this.T(), target.NewRuntime("runtime"))

	target.Finalize()
}
