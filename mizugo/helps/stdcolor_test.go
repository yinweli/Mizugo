package helps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestStdColor(t *testing.T) {
	suite.Run(t, new(SuiteStdColor))
}

type SuiteStdColor struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteStdColor) SetupSuite() {
	this.Env = testdata.EnvSetup("test-helps-stdColor")
}

func (this *SuiteStdColor) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteStdColor) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteStdColor) TestStdColor() {
	target := NewStdColor(os.Stdout, os.Stderr)
	assert.NotNil(this.T(), target)
	assert.NotNil(this.T(), target.Outf("test outf %v\n", 100))
	assert.NotNil(this.T(), target.Outln("test outln"))
	assert.NotNil(this.T(), target.Errf("test errf %v\n", 100))
	assert.NotNil(this.T(), target.Errln("test errln"))
	assert.NotNil(this.T(), target.GetStdout())
	assert.NotNil(this.T(), target.GetStderr())
	_, _ = target.GetStdout().Write([]byte("test stdout\n"))
	_, _ = target.GetStderr().Write([]byte("test stderr\n"))
}
