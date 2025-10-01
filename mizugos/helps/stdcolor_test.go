package helps

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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
	this.NotNil(target)
	target = NewStdColor(nil, os.Stderr)
	this.NotNil(target)
	target = NewStdColor(os.Stdout, nil)
	this.NotNil(target)
}

func (this *SuiteStdColor) TestOut() {
	target := NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Out("test out %v", 100))
	this.False(target.Failed())

	target = NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Outf("test outf %v\n", 100))
	this.False(target.Failed())

	target = NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Outln("test outln"))
	this.False(target.Failed())
}

func (this *SuiteStdColor) TestErr() {
	target := NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Err("test err %v", 100))
	this.True(target.Failed())

	target = NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Errf("test errf %v\n", 100))
	this.True(target.Failed())

	target = NewStdColor(os.Stdout, os.Stderr)
	this.NotNil(target.Errln("test errln"))
	this.True(target.Failed())
}

func (this *SuiteStdColor) TestStdout() {
	target := NewStdColor(os.Stdout, os.Stderr).GetStdout()
	this.NotNil(target)
	_, _ = target.Write([]byte("test stdout\n"))
}

func (this *SuiteStdColor) TestStderr() {
	target := NewStdColor(os.Stdout, os.Stderr).GetStderr()
	this.NotNil(target)
	_, _ = target.Write([]byte("test stderr\n"))
}
