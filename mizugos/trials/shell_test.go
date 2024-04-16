package trials

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
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
	assert.Equal(this.T(), 0, Shell(nil, proc, proc, work))
	assert.Equal(this.T(), 2, count)
}

func (this *SuiteShell) TestPrepareRestore() {
	work := filepath.Join(Root(), "prepare-restore")
	test := filepath.Join(Root(), "prepare-restore", testdata.TrialFileName)
	catalog := Prepare(work, testdata.TrialDir)
	assert.True(this.T(), FileExist(test))
	Restore(catalog)
	assert.False(this.T(), FileExist(test))
}

func (this *SuiteShell) TestRoot() {
	assert.Contains(this.T(), Root(), filepath.Clean("mizugos/trials"))
}
