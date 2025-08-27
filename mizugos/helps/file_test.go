package helps

import (
	"io/fs"
	"os"
	"testing"

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

func (this *SuiteFile) TestFileExist() {
	this.Nil(FileWrite("file/exist.txt", []byte("ok")))
	this.True(FileExist("file/exist.txt"))
	this.False(FileExist(testdata.Unknown))
}

func (this *SuiteFile) TestFileCompare() {
	this.Nil(FileWrite("file/compare.txt", []byte("ok")))
	this.True(FileCompare("file/compare.txt", []byte("ok")))
	this.False(FileCompare("file/compare.txt", []byte("no")))
	this.False(FileCompare(testdata.Unknown, []byte("ok")))
}

func (this *SuiteFile) TestFileWrite() {
	_ = os.WriteFile("parent", []byte("x"), fs.ModePerm)
	this.Nil(FileWrite("file/write.txt", []byte("ok")))
	this.NotNil(FileWrite("parent/write.txt", []byte("ok")))
	this.NotNil(FileWrite("bad\x00name.txt", []byte("ok")))
}
