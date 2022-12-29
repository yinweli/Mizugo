package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMetricsmgr(t *testing.T) {
	suite.Run(t, new(SuiteMetricsmgr))
}

type SuiteMetricsmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMetricsmgr) SetupSuite() {
	this.Change("test-metrics-metricsmgr")
}

func (this *SuiteMetricsmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMetricsmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteMetricsmgr) TestInitialize() {
	target := NewMetricsmgr()
	target.Initialize(8080, &Auth{
		Username: "username",
		Password: "password",
	})
	target.Finalize()
}

func (this *SuiteMetricsmgr) TestNew() {
	target := NewMetricsmgr()
	target.Initialize(8080, nil)

	assert.NotNil(this.T(), target.NewInt("int"))
	assert.NotNil(this.T(), target.NewFloat("float"))
	assert.NotNil(this.T(), target.NewString("string"))
	assert.NotNil(this.T(), target.NewMap("map"))
	assert.NotNil(this.T(), target.NewRuntime("runtime"))

	target.Finalize()
}
