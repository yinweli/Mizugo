package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestProcmgr(t *testing.T) {
	suite.Run(t, new(SuiteProcmgr))
}

type SuiteProcmgr struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteProcmgr) SetupSuite() {
	this.Env = testdata.EnvSetup("test-procs-procmgr")
}

func (this *SuiteProcmgr) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteProcmgr) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteProcmgr) TestNewProcmgr() {
	assert.NotNil(this.T(), NewProcmgr())
}

func (this *SuiteProcmgr) TestProcmgr() {
	target := NewProcmgr()
	messageID := MessageID(1)
	target.Add(messageID, func(_ any) {})
	assert.NotNil(this.T(), target.Get(messageID))
	target.Del(messageID)
	assert.Nil(this.T(), target.Get(messageID))
}
