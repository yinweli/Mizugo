package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/yinweli/Mizugo/mizugos/utils"
	"github.com/yinweli/Mizugo/testdata"
)

func TestProtoDes(t *testing.T) {
	suite.Run(t, new(SuiteProtoDes))
}

type SuiteProtoDes struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteProtoDes) SetupSuite() {
	this.Change("test-procs-protodes")
}

func (this *SuiteProtoDes) TearDownSuite() {
	this.Restore()
}

func (this *SuiteProtoDes) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteProtoDes) TestNewProtoDes() {
	assert.NotNil(this.T(), NewProtoDes([]byte{}))
}

func (this *SuiteProtoDes) TestEncodeDecode() {
	key := utils.DesKeyRand()
	target := NewProtoDes(key)
	msg := protoDesTestMsg("test encode/decode message")

	packet, err := target.Encode(msg)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), packet)

	result, err := target.Decode(packet)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), result)
	assert.True(this.T(), proto.Equal(msg, result.(*ProtoDesMsg)))

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown packet data"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProtoDes) TestProcess() {
	key := utils.DesKeyRand()
	target := NewProtoDes(key)
	msg := protoDesTestMsg("test process message")

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

func BenchmarkProtoDesEncode(b *testing.B) {
	key := utils.DesKeyRand()
	target := NewProtoDes(key)
	msg := protoDesTestMsg("benchmark encode message")

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(msg)
	} // for
}

func BenchmarkProtoDesDecode(b *testing.B) {
	key := utils.DesKeyRand()
	target := NewProtoDes(key)
	msg := protoDesTestMsg("benchmark decode message")
	packet, _ := target.Encode(msg)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(packet)
	} // for
}

func protoDesTestMsg(message string) *ProtoDesMsg {
	msg, _ := anypb.New(&ProtoDesMsgTest{
		Message: message,
	})
	return &ProtoDesMsg{
		MessageID: 1,
		Message:   msg,
	}
}
