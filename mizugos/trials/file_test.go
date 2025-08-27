package trials

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestFile(t *testing.T) {
	suite.Run(t, new(SuiteFile))
}

type SuiteFile struct {
	suite.Suite
}

func (this *SuiteFile) TestFileExist() {
	assert.True(this.T(), FileExist(testdata.TrialFile))
	assert.False(this.T(), FileExist(testdata.Unknown))
}

func (this *SuiteFile) TestFileCompare() {
	assert.True(this.T(), FileCompare(testdata.TrialFile, []byte("0")))
	assert.False(this.T(), FileCompare(testdata.TrialFile, []byte("1")))
	assert.False(this.T(), FileCompare(testdata.Unknown, []byte("0")))
}
