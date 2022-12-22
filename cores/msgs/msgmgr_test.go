package msgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestMsgmgr(t *testing.T) {
	suite.Run(t, new(SuiteMsgmgr))
}

type SuiteMsgmgr struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteMsgmgr) SetupSuite() {
	this.Change("test-msgs-msgmgr")
}

func (this *SuiteMsgmgr) TearDownSuite() {
	this.Restore()
}

func (this *SuiteMsgmgr) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteMsgmgr) TestNewMsgmgr() {
	assert.NotNil(this.T(), NewMsgmgr())
}

func (this *SuiteMsgmgr) TestMsgmgr() {
	target := NewMsgmgr()
	messageID := MessageID(1)
	target.Add(messageID, func(messageID MessageID, message any) {
		// do nothing
	})
	assert.NotNil(this.T(), target.Get(messageID))
	target.Del(messageID)
	assert.Nil(this.T(), target.Get(messageID))
}

func (this *SuiteMsgmgr) TestCast() {
	msg := &msgTester1{}

	result, err := Cast[msgTester1](msg)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), result)

	_, err = Cast[msgTester2](msg)
	assert.NotNil(this.T(), err)

	_, err = Cast[msgTester1](nil)
	assert.NotNil(this.T(), err)

	_, err = Cast[msgTester2](nil)
	assert.NotNil(this.T(), err)
}

type msgTester1 struct {
}

type msgTester2 struct {
}
