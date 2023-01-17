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
	testdata.TestEnv
	testdata.TestLeak
}

func (this *SuiteProcmgr) SetupSuite() {
	this.Change("test-procs-procmgr")
}

func (this *SuiteProcmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteProcmgr) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteProcmgr) TestNewProcmgr() {
	assert.NotNil(this.T(), newProcmgr())
}

func (this *SuiteProcmgr) TestProcmgr() {
	target := newProcmgr()
	messageID := MessageID(1)
	target.Add(messageID, func(message any) {
		// do nothing
	})
	assert.NotNil(this.T(), target.Get(messageID))
	target.Del(messageID)
	assert.Nil(this.T(), target.Get(messageID))
}
