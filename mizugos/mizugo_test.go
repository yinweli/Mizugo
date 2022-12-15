package mizugos

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMizugo(t *testing.T) {
	suite.Run(t, new(SuiteMizugo))
}

type SuiteMizugo struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMizugo) SetupSuite() {
	this.Change("test-mizugo")
}

func (this *SuiteMizugo) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMizugo) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteMizugo) TestStart() {
	/*
		Initialize(&logs.EmptyLogger{})
		assert.NotNil(this.T(), Logger())
		assert.NotNil(this.T(), Netmgr())
		assert.NotNil(this.T(), Entitymgr())
		assert.NotNil(this.T(), Tagmgr())
	*/
}
