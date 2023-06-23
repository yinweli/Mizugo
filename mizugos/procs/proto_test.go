package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestProto(t *testing.T) {
	suite.Run(t, new(SuiteProto))
}

type SuiteProto struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteProto) SetupSuite() {
	this.Env = testdata.EnvSetup("test-procs-proto")
}

func (this *SuiteProto) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteProto) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteProto) TestEncode() {
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	target := NewProto()
	assert.NotNil(this.T(), target)
	input, _ := ProtoMarshal(messageID, message)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Encode("!?")
	assert.NotNil(this.T(), err)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(input, decode.(*msgs.ProtoMsg)))

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestEncodeBase64() {
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	target := NewProto().Base64(true)
	input, _ := ProtoMarshal(messageID, message)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Encode("!?")
	assert.NotNil(this.T(), err)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(input, decode.(*msgs.ProtoMsg)))

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestEncodeDesCBC() {
	key := cryptos.RandDesKeyString()
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	target := NewProto().DesCBC(true, key, key) // 這裡偷懶把key跟iv都設為key
	input, _ := ProtoMarshal(messageID, message)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Encode("!?")
	assert.NotNil(this.T(), err)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(input, decode.(*msgs.ProtoMsg)))

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestEncodeAll() {
	key := cryptos.RandDesKeyString()
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	target := NewProto().Base64(true).DesCBC(true, key, key) // 這裡偷懶把key跟iv都設為key
	input, _ := ProtoMarshal(messageID, message)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Encode("!?")
	assert.NotNil(this.T(), err)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(input, decode.(*msgs.ProtoMsg)))

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestProcess() {
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	target := NewProto()
	input, _ := ProtoMarshal(messageID, message)

	valid := false
	target.Add(messageID, func(message any) {
		valid = proto.Equal(input, message.(*msgs.ProtoMsg))
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), valid)

	input, _ = ProtoMarshal(0, message)
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuiteProto) TestMarshal() {
	messageID := MessageID(1)
	message := &msgs.ProtoTest{
		Data: "proto test",
	}
	output1, err := ProtoMarshal(messageID, message)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)

	_, err = ProtoMarshal(messageID, nil)
	assert.NotNil(this.T(), err)

	messageID, output2, err := ProtoUnmarshal[msgs.ProtoTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), messageID, messageID)
	assert.True(this.T(), proto.Equal(message, output2))

	_, _, err = ProtoUnmarshal[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)

	_, _, err = ProtoUnmarshal[msgs.ProtoTest]("!?")
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

func BenchmarkProtoEncodeBase64(b *testing.B) {
	target := NewProto().Base64(true)
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkProtoEncodeDesCBC(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewProto().DesCBC(true, key, key)
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkProtoEncodeAll(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewProto().Base64(true).DesCBC(true, key, key)
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

func BenchmarkProtoDecodeBase64(b *testing.B) {
	target := NewProto().Base64(true)
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoDecodeDesCBC(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewProto().DesCBC(true, key, key)
	input, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoDecodeAll(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewProto().Base64(true).DesCBC(true, key, key)
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
