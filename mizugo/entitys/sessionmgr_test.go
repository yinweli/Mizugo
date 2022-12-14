package entitys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/mizugo/nets"
	"github.com/yinweli/Mizugo/testdata"
)

func TestSessionmgr(t *testing.T) {
	suite.Run(t, new(SuiteSessionmgr))
}

type SuiteSessionmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteSessionmgr) SetupSuite() {
	this.Change("test-entitys-sessionmgr")
}

func (this *SuiteSessionmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSessionmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteSessionmgr) TestNewSessionmgr() {
	assert.NotNil(this.T(), NewSessionmgr())
}

func (this *SuiteSessionmgr) TestSetGet() {
	session := nets.NewTCPSession(nil)
	target := NewSessionmgr()
	target.Set(session)
	assert.Equal(this.T(), session, target.Get())
}
