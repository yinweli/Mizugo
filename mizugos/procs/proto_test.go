package procs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
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

func (this *SuiteProto) TestEncode() {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	this.NotNil(target)
	encode, err := target.Encode(marshal)
	this.Nil(err)
	this.NotNil(encode)
	_, err = target.Encode(nil)
	this.NotNil(err)
	_, err = target.Encode(testdata.Unknown)
	this.NotNil(err)
}

func (this *SuiteProto) TestDecode() {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	encode, _ := target.Encode(marshal)
	decode, err := target.Decode(encode)
	this.Nil(err)
	this.NotNil(decode)
	this.True(proto.Equal(marshal, decode.(*msgs.Proto)))
	_, err = target.Decode(nil)
	this.NotNil(err)
	_, err = target.Decode(testdata.Unknown)
	this.NotNil(err)
	_, err = target.Decode([]byte(testdata.Unknown))
	this.NotNil(err)
}

func (this *SuiteProto) TestProcess() {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	valid := false
	target := NewProto()
	target.Add(1, func(message any) {
		valid = proto.Equal(marshal, message.(*msgs.Proto))
	})
	this.Nil(target.Process(marshal))
	this.True(valid)
	marshal, _ = ProtoMarshal(2, &msgs.ProtoTest{})
	this.NotNil(target.Process(marshal))
	this.NotNil(target.Process(nil))
	this.NotNil(target.Process(testdata.Unknown))
}

func (this *SuiteProto) TestMarshal() {
	marshal, err := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	this.Nil(err)
	this.NotNil(marshal)
	_, err = ProtoMarshal(1, nil)
	this.NotNil(err)
}

func (this *SuiteProto) TestUnmarshal() {
	message := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}
	marshal, _ := ProtoMarshal(1, message)
	messageID, payload, err := ProtoUnmarshal[*msgs.ProtoTest](marshal)
	this.Nil(err)
	this.Equal(int32(1), messageID)
	this.True(proto.Equal(message, payload))
	_, _, err = ProtoUnmarshal[*msgs.ProtoTest](&msgs.ProtoTest{})
	this.NotNil(err)
	_, _, err = ProtoUnmarshal[*msgs.ProtoTest](nil)
	this.NotNil(err)
}

func BenchmarkProtoEncode(b *testing.B) {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(marshal)
	} // for
}

func BenchmarkProtoDecode(b *testing.B) {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	encode, _ := target.Encode(marshal)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoMarshal(b *testing.B) {
	message := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = ProtoMarshal(1, message)
	} // for
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	marshal, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = ProtoUnmarshal[*msgs.ProtoTest](marshal)
	} // for
}
