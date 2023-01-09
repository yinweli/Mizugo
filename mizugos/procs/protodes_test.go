package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"
	"google.golang.org/protobuf/proto"

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
	assert.NotNil(this.T(), NewProtoDes())
}

func (this *SuiteProtoDes) TestEncodeDecode() {
	target := NewProtoDes()
	target.Key(utils.RandDesKey())
	input, err := ProtoDesMarshal(1, &ProtoDesMsgTest{
		Message: "test encode/decode",
	})
	assert.Nil(this.T(), err)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(input, decode.(*ProtoDesMsg)))

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode/decode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProtoDes) TestProcess() {
	target := NewProtoDes()
	target.Key(utils.RandDesKey())
	input, err := ProtoDesMarshal(1, &ProtoDesMsgTest{
		Message: "test process",
	})
	assert.Nil(this.T(), err)

	valid := false
	target.Add(input.MessageID, func(message any) {
		valid = proto.Equal(input, message.(*ProtoDesMsg))
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), valid)

	input.MessageID = 2
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuiteProtoDes) TestMarshal() {
	inputID := MessageID(1)
	input := &ProtoDesMsgTest{
		Message: "test marshal/unmarshal",
	}

	message, err := ProtoDesMarshal(inputID, input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), message)

	outputID, output, err := ProtoDesUnmarshal[*ProtoDesMsgTest](message)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), inputID, outputID)
	assert.True(this.T(), proto.Equal(input, output))

	_, _, err = ProtoDesUnmarshal[*ProtoDesMsgTest](nil)
	assert.NotNil(this.T(), err)

	_, _, err = ProtoDesUnmarshal[*ProtoDesMsg](message)
	assert.NotNil(this.T(), err)
}

func BenchmarkProtoDesEncode(b *testing.B) {
	target := NewProtoDes()
	target.Key(utils.RandDesKey())
	input, _ := ProtoDesMarshal(1, &ProtoDesMsgTest{
		Message: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkProtoDesDecode(b *testing.B) {
	target := NewProtoDes()
	target.Key(utils.RandDesKey())
	input, _ := ProtoDesMarshal(1, &ProtoDesMsgTest{
		Message: "benchmark encode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoDesMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ProtoDesMarshal(1, &ProtoDesMsgTest{
			Message: "benchmark marshal",
		})
	} // for
}

func BenchmarkProtoDesUnmarshal(b *testing.B) {
	input, _ := ProtoDesMarshal(1, &ProtoDesMsgTest{
		Message: "benchmark unmarshal",
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = ProtoDesUnmarshal[*ProtoDesMsgTest](input)
	} // for
}
