package msgs

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestStringMsg(t *testing.T) {
	suite.Run(t, new(SuiteStringMsg))
}

type SuiteStringMsg struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteStringMsg) SetupSuite() {
	this.Change("test-msgs-stringmsg")
}

func (this *SuiteStringMsg) TearDownSuite() {
	this.Restore()
}

func (this *SuiteStringMsg) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteStringMsg) TestNewStringMsg() {
	assert.NotNil(this.T(), NewStringMsg())
}

func (this *SuiteStringMsg) TestEncodeDecode() {
	target := NewStringMsg()
	msg := &StringMessage{
		MessageID: 1,
		Message:   "test encode/decode message",
	}
	packet, err := target.Encode(msg)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), packet)
	result, err := target.Decode(packet)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), result)
	assert.Equal(this.T(), msg, result)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown packet data"))
	assert.NotNil(this.T(), err)

	msg.Sum = "unknown md5 string"
	bytes, _ := json.Marshal(msg)
	packet = []byte(base64.StdEncoding.EncodeToString(bytes))
	_, err = target.Decode(packet)
	assert.NotNil(this.T(), err)
}

func (this *SuiteStringMsg) TestProcess() {
	target := NewStringMsg()
	msg := &StringMessage{
		MessageID: 1,
		Message:   "test process message",
	}
	valid := false
	target.Add(msg.MessageID, func(messageID MessageID, message any) {
		assert.Equal(this.T(), msg, message)
		valid = true
	})
	assert.Nil(this.T(), target.Process(msg))
	assert.True(this.T(), valid)

	msg.MessageID = 2
	assert.NotNil(this.T(), target.Process(msg))

	assert.NotNil(this.T(), target.Process(nil))
}

func BenchmarkStringMsgEncode(b *testing.B) {
	target := NewStringMsg()
	msg := &StringMessage{
		MessageID: 1,
		Message:   "benchmark encode message",
	}

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(msg)
	} // for
}

func BenchmarkStringMsgDecode(b *testing.B) {
	target := NewStringMsg()
	msg := &StringMessage{
		MessageID: 1,
		Message:   "benchmark decode message",
	}
	packet, _ := target.Encode(msg)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(packet)
	} // for
}
