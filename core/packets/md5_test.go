package packets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMD5(t *testing.T) {
	suite.Run(t, new(SuiteMD5))
}

type SuiteMD5 struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMD5) SetupSuite() {
	this.Change("test-packets-md5")
}

func (this *SuiteMD5) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMD5) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteMD5) TestMD5String() {
	assert.NotNil(this.T(), MD5String([]byte("12345")))
}
