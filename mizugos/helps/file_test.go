package helps

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestFile(t *testing.T) {
	suite.Run(t, new(SuiteFile))
}

type SuiteFile struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteFile) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-helps-file"))
}

func (this *SuiteFile) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteFile) TestWriteFile() {
	path := filepath.Join("path", "test.file")
	data := []byte("this is a string")

	assert.Nil(this.T(), WriteFile(path, data))
	assert.True(this.T(), trials.FileCompare(path, data))
}
