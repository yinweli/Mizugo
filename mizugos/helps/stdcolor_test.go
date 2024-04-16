package helps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestStdColor(t *testing.T) {
	suite.Run(t, new(SuiteStdColor))
}

type SuiteStdColor struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteStdColor) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-stdColor"))
}

func (this *SuiteStdColor) TearDownSuite() {
	trials.Restore(this.Catalog)
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
