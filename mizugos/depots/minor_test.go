package depots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/contexts"
	"github.com/yinweli/Mizugo/testdata"
)

func TestMinor(t *testing.T) {
	suite.Run(t, new(SuiteMinor))
}

type SuiteMinor struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteMinor) SetupSuite() {
	this.Change("test-depots-minor")
}

func (this *SuiteMinor) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMinor) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteMinor) TestNewMinor() {
	target, err := newMinor(contexts.Ctx(), "mongodb://127.0.0.1:27017/")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target)
}

func (this *SuiteMinor) TestMinor() {
	target, err := newMinor(contexts.Ctx(), "mongodb://127.0.0.1:27017/")
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), target.Runner("database", "table"))
	assert.NotNil(this.T(), target.Client())
	target.stop(contexts.Ctx())
	assert.Nil(this.T(), target.Runner("database", "table"))
	assert.Nil(this.T(), target.Client())

	_, err = newMinor(contexts.Ctx(), "mongodb://127.0.0.1:10001/?timeoutMS=1000")
	assert.NotNil(this.T(), err)
}
