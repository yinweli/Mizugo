package redmos

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMX(t *testing.T) {
	suite.Run(t, new(SuiteMX))
}

type SuiteMX struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteMX) SetupSuite() {
	this.Env = testdata.EnvSetup("test-redmos-mx")
}

func (this *SuiteMX) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteMX) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteMX) TestSave() {
	target := &Save{}
	target.Set()
	assert.True(this.T(), target.Get())
}
