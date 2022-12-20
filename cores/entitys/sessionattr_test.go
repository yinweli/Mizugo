package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/cores/nets"
	"github.com/yinweli/Mizugo/testdata"
)

func TestSessionAttr(t *testing.T) {
	suite.Run(t, new(SuiteSessionAttr))
}

type SuiteSessionAttr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteSessionAttr) SetupSuite() {
	this.Change("test-entitys-sessionattr")
}

func (this *SuiteSessionAttr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSessionAttr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteSessionAttr) TestSetGet() {
	session := nets.NewTCPSession(nil)
	target := SessionAttr{}
	target.Set(session)
	assert.Equal(this.T(), session, target.Get())
}
