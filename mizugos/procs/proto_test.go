package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
)

func TestProto(t *testing.T) {
	suite.Run(t, new(SuiteProto))
}

type SuiteProto struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteProto) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-procs-proto"))
}

func (this *SuiteProto) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteProto) TestProto() {
	target := NewProto()
	assert.NotNil(this.T(), target)
	output, _ := ProtoMarshal(MessageID(1), &msgs.ProtoTest{})
	encode, err := target.Encode(output)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)
	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Encode("!?")
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestDecode() {
	target := NewProto()
	output, _ := ProtoMarshal(MessageID(1), &msgs.ProtoTest{})
	encode, _ := target.Encode(output)
	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(output, decode.(*msgs.ProtoMsg)))
	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = target.Decode([]byte(testdata.Unknown))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestProcess() {
	target := NewProto()
	output, _ := ProtoMarshal(MessageID(1), &msgs.ProtoTest{})
	valid := false
	target.Add(MessageID(1), func(message any) {
		valid = proto.Equal(output, message.(*msgs.ProtoMsg))
	})
	assert.Nil(this.T(), target.Process(output))
	assert.True(this.T(), valid)
	output2, _ := ProtoMarshal(MessageID(2), &msgs.ProtoTest{})
	assert.NotNil(this.T(), target.Process(output2))
	assert.NotNil(this.T(), target.Process(nil))
	assert.NotNil(this.T(), target.Process(testdata.Unknown))
}

func (this *SuiteProto) TestMarshal() {
	output, err := ProtoMarshal(MessageID(1), &msgs.ProtoTest{})
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output)
	_, err = ProtoMarshal(MessageID(1), nil)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestUnmarshal() {
	message := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}
	output1, _ := ProtoMarshal(MessageID(1), message)
	messageID, output2, err := ProtoUnmarshal[msgs.ProtoTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), MessageID(1), messageID)
	assert.True(this.T(), proto.Equal(message, output2))
	_, _, err = ProtoUnmarshal[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)
	_, _, err = ProtoUnmarshal[int](output1)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestAny() {
	message := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}
	output1, _ := ProtoMarshal(MessageID(1), message)
	output2, err := ProtoAny[msgs.ProtoTest](output1.Message)
	assert.Nil(this.T(), err)
	assert.True(this.T(), proto.Equal(message, output2))
	_, err = ProtoAny[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)
	_, err = ProtoAny[int](output1.Message)
	assert.NotNil(this.T(), err)
}

func BenchmarkProtoEncode(b *testing.B) {
	target := NewProto()
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkProtoDecode(b *testing.B) {
	target := NewProto()
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoMarshal(b *testing.B) {
	input := &msgs.ProtoTest{
		Data: "benchmark marshal",
	}

	for i := 0; i < b.N; i++ {
		_, _ = ProtoMarshal(1, input)
	} // for
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark unmarshal",
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = ProtoUnmarshal[msgs.ProtoTest](input)
	} // for
}
