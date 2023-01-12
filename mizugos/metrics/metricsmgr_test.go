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
}

func (this *SuiteMetricsmgr) SetupSuite() {
	this.Change("test-metrics-metricsmgr")
}

func (this *SuiteMetricsmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMetricsmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMetricsmgr) TestInitialize() {
	target := NewMetricsmgr()
	target.Initialize(8080)
	target.Finalize()
}

func (this *SuiteMetricsmgr) TestNew() {
	target := NewMetricsmgr()
	target.Initialize(8080)

	assert.NotNil(this.T(), target.NewInt("int"))
	assert.NotNil(this.T(), target.NewFloat("float"))
	assert.NotNil(this.T(), target.NewString("string"))
	assert.NotNil(this.T(), target.NewMap("map"))
	assert.NotNil(this.T(), target.NewRuntime("runtime"))

	target.Finalize()
}
