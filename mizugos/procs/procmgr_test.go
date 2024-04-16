package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestProcmgr(t *testing.T) {
	suite.Run(t, new(SuiteProcmgr))
}

type SuiteProcmgr struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteProcmgr) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-procs-procmgr"))
}

func (this *SuiteProcmgr) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteProcmgr) TestNewProcmgr() {
	target := NewProcmgr()
	assert.NotNil(this.T(), target)
	messageID := MessageID(1)
	target.Add(messageID, func(_ any) {})
	assert.NotNil(this.T(), target.Get(messageID))
	target.Del(messageID)
	assert.Nil(this.T(), target.Get(messageID))
}
