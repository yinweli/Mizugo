package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	assert.NotNil(this.T(), target)
	encode, err := target.Encode(output)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)
	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestDecode() {
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	encode, _ := target.Encode(output)
	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(output, decode.(*msgs.Proto)))
	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = target.Decode([]byte(testdata.Unknown))
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestProcess() {
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	valid := false
	target := NewProto()
	target.Add(1, func(message any) {
		valid = proto.Equal(output, message.(*msgs.Proto))
	})
	assert.Nil(this.T(), target.Process(output))
	assert.True(this.T(), valid)
	output, _ = ProtoMarshal(2, &msgs.ProtoTest{})
	assert.NotNil(this.T(), target.Process(output))
	assert.NotNil(this.T(), target.Process(nil))
	assert.NotNil(this.T(), target.Process(testdata.Unknown))
}

func (this *SuiteProto) TestMarshal() {
	output, err := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output)
	_, err = ProtoMarshal(1, nil)
	assert.NotNil(this.T(), err)
}

func (this *SuiteProto) TestUnmarshal() {
	object := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}
	output1, _ := ProtoMarshal(1, object)
	messageID, output2, err := ProtoUnmarshal[msgs.ProtoTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), int32(1), messageID)
	assert.True(this.T(), proto.Equal(object, output2))
	_, _, err = ProtoUnmarshal[msgs.ProtoTest](nil)
	assert.NotNil(this.T(), err)
	_, _, err = ProtoUnmarshal[int](output1)
	assert.NotNil(this.T(), err)
}

func BenchmarkProtoEncode(b *testing.B) {
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(output)
	} // for
}

func BenchmarkProtoDecode(b *testing.B) {
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})
	target := NewProto()
	encode, _ := target.Encode(output)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkProtoMarshal(b *testing.B) {
	object := &msgs.ProtoTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = ProtoMarshal(1, object)
	} // for
}

func BenchmarkProtoUnmarshal(b *testing.B) {
	output, _ := ProtoMarshal(1, &msgs.ProtoTest{
		Data: testdata.Unknown,
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = ProtoUnmarshal[msgs.ProtoTest](output)
	} // for
}
