package trials

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestShell(t *testing.T) {
	suite.Run(t, new(SuiteShell))
}

type SuiteShell struct {
	suite.Suite
}

func (this *SuiteShell) TestShell() {
	work := filepath.Join(Root(), "shell")
	count := 0
	proc := func() {
		count++
	}
	this.Equal(0, Shell(nil, proc, proc, work))
	this.Equal(2, count)
}

func (this *SuiteShell) TestPrepareRestore() {
	work := filepath.Join(Root(), "prepare-restore")
	path := filepath.Join(Root(), "prepare-restore", testdata.TrialFileName)
	catalog := Prepare(work, testdata.TrialDir)
	this.FileExists(path)
	Restore(catalog)
	this.NoFileExists(path)
}

func (this *SuiteShell) TestRoot() {
	this.Contains(Root(), filepath.Clean("mizugos/trials"))
}
