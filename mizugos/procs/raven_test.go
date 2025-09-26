package procs

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestRaven(t *testing.T) {
	suite.Run(t, new(SuiteRaven))
}

type SuiteRaven struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteRaven) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-procs-raven"))
}

func (this *SuiteRaven) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteRaven) TestEncode() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message, message)
	target := NewRaven()
	this.NotNil(target)
	encode, err := target.Encode(builder)
	this.Nil(err)
	this.NotNil(encode)
	this.IsType([]byte{}, encode)
	_, err = target.Encode(nil)
	this.NotNil(err)
	_, err = target.Encode(testdata.Unknown)
	this.NotNil(err)
}

func (this *SuiteRaven) TestDecode() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)
	marshal, _ := proto.Marshal(builder.(proto.Message))
	target := NewRaven()
	decode, err := target.Decode(marshal)
	this.Nil(err)
	this.NotNil(decode)
	this.True(proto.Equal(builder.(proto.Message), decode.(proto.Message)))
	_, err = target.Decode(nil)
	this.NotNil(err)
	_, err = target.Decode(testdata.Unknown)
	this.NotNil(err)
	_, err = target.Decode([]byte(testdata.Unknown))
	this.NotNil(err)
}

func (this *SuiteRaven) TestProcess() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)
	valid := false
	target := NewRaven()
	target.Add(1, func(message any) {
		valid = proto.Equal(builder.(proto.Message), message.(proto.Message))
	})
	this.Nil(target.Process(builder))
	this.True(valid)
	builder, _ = RavenSBuilder(2, message, message)
	this.NotNil(target.Process(builder))
	this.NotNil(target.Process(nil))
	this.NotNil(target.Process(testdata.Unknown))
}

func (this *SuiteRaven) TestClientEncode() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)
	target := NewRavenClient()
	this.NotNil(target)
	encode, err := target.Encode(builder)
	this.Nil(err)
	this.NotNil(encode)
	this.IsType([]byte{}, encode)
	_, err = target.Encode(nil)
	this.NotNil(err)
	_, err = target.Encode(testdata.Unknown)
	this.NotNil(err)
}

func (this *SuiteRaven) TestClientDecode() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message)
	marshal, _ := proto.Marshal(builder.(proto.Message))
	target := NewRavenClient()
	decode, err := target.Decode(marshal)
	this.Nil(err)
	this.NotNil(decode)
	this.True(proto.Equal(builder.(proto.Message), decode.(proto.Message)))
	_, err = target.Decode(nil)
	this.NotNil(err)
	_, err = target.Decode(testdata.Unknown)
	this.NotNil(err)
	_, err = target.Decode([]byte(testdata.Unknown))
	this.NotNil(err)
}

func (this *SuiteRaven) TestClientProcess() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message)
	valid := false
	target := NewRavenClient()
	target.Add(1, func(message any) {
		valid = proto.Equal(builder.(proto.Message), message.(proto.Message))
	})
	this.Nil(target.Process(builder))
	this.True(valid)
	builder, _ = RavenCBuilder(2, 1, message, message, message)
	this.NotNil(target.Process(builder))
	this.NotNil(target.Process(nil))
	this.NotNil(target.Process(testdata.Unknown))
}

func (this *SuiteRaven) TestRavenS() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, err := RavenSBuilder(1, message, message)
	this.Nil(err)
	this.NotNil(builder)
	_, err = RavenSBuilder(1, nil, message)
	this.NotNil(err)
	_, err = RavenSBuilder(1, message, nil)
	this.NotNil(err)
	parser, err := RavenSParser[*msgs.RavenTest, *msgs.RavenTest](builder)
	this.Nil(err)
	this.NotNil(parser)
	this.NotZero(parser.Size())
	this.NotNil(parser.Detail())
	this.Equal(int32(1), parser.MessageID)
	this.True(proto.Equal(message, parser.Header))
	this.True(proto.Equal(message, parser.Request))
	_, err = RavenSParser[*msgs.RavenTest, *msgs.RavenTest](nil)
	this.NotNil(err)
	_, err = RavenSParser[*msgs.RavenTest, *msgs.RavenTest](testdata.Unknown)
	this.NotNil(err)
	_, err = RavenSParser[*msgs.ProtoTest, *msgs.RavenTest](builder)
	this.NotNil(err)
	_, err = RavenSParser[*msgs.RavenTest, *msgs.ProtoTest](builder)
	this.NotNil(err)
}

func (this *SuiteRaven) TestRavenC() {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, err := RavenCBuilder(1, 1, message, message, message, message)
	this.Nil(err)
	this.NotNil(builder)
	_, err = RavenCBuilder(1, 1, nil, message)
	this.NotNil(err)
	_, err = RavenCBuilder(1, 1, message, nil)
	this.NotNil(err)
	_, err = RavenCBuilder(1, 1, message, message, nil)
	this.NotNil(err)
	parser, err := RavenCParser[*msgs.RavenTest, *msgs.RavenTest](builder)
	this.Nil(err)
	this.NotNil(parser)
	this.NotZero(parser.Size())
	this.NotNil(parser.Detail())
	this.Equal(int32(1), parser.MessageID)
	this.Equal(int32(1), parser.ErrID)
	this.True(proto.Equal(message, parser.Header))
	this.True(proto.Equal(message, parser.Request))
	this.True(proto.Equal(message, parser.Respond[0]))
	this.True(proto.Equal(message, parser.Respond[1]))
	_, err = RavenCParser[*msgs.RavenTest, *msgs.RavenTest](nil)
	this.NotNil(err)
	_, err = RavenCParser[*msgs.RavenTest, *msgs.RavenTest](testdata.Unknown)
	this.NotNil(err)
	_, err = RavenCParser[*msgs.ProtoTest, *msgs.RavenTest](builder)
	this.NotNil(err)
	_, err = RavenCParser[*msgs.RavenTest, *msgs.ProtoTest](builder)
	this.NotNil(err)

	this.True(RavenIsMessageID(builder, 1))
	this.True(RavenIsErrID(builder, 1))
	this.True(proto.Equal(message, RavenHeader[*msgs.RavenTest](builder)))
	this.Nil(RavenHeader[*msgs.RavenTest](nil))
	this.Nil(RavenHeader[*msgs.RavenTest](testdata.Unknown))
	this.Nil(RavenHeader[*msgs.ProtoTest](builder))
	this.True(proto.Equal(message, RavenRequest[*msgs.RavenTest](builder)))
	this.Nil(RavenRequest[*msgs.RavenTest](nil))
	this.Nil(RavenRequest[*msgs.RavenTest](testdata.Unknown))
	this.Nil(RavenRequest[*msgs.ProtoTest](builder))
	this.True(proto.Equal(message, RavenRespondAt[*msgs.RavenTest](builder, 0)))
	this.Nil(RavenRespondAt[*msgs.RavenTest](nil, 1))
	this.Nil(RavenRespondAt[*msgs.RavenTest](testdata.Unknown, 1))
	this.Nil(RavenRespondAt[*msgs.RavenTest](builder, -1))
	this.Nil(RavenRespondAt[*msgs.ProtoTest](builder, 1))
	this.True(proto.Equal(message, RavenRespondFind[*msgs.RavenTest](builder)))
	this.Nil(RavenRespondFind[*msgs.RavenTest](nil))
	this.Nil(RavenRespondFind[*msgs.RavenTest](testdata.Unknown))
	this.Nil(RavenRespondFind[*msgs.ProtoTest](builder))
}

func BenchmarkRavenEncode(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message)
	target := NewRaven()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(builder)
	} // for
}

func BenchmarkRavenDecode(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)
	target := NewRaven()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(builder)
	} // for
}

func BenchmarkRavenClientEncode(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)
	target := NewRavenClient()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(builder)
	} // for
}

func BenchmarkRavenClientDecode(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message)
	target := NewRavenClient()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(builder)
	} // for
}

func BenchmarkRavenSBuilder(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = RavenSBuilder(1, message, message)
	} // for
}

func BenchmarkRavenSParser(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenSBuilder(1, message, message)

	for i := 0; i < b.N; i++ {
		_, _ = RavenSParser[*msgs.RavenTest, *msgs.RavenTest](builder)
	} // for
}

func BenchmarkRavenCBuilder(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = RavenCBuilder(1, 1, message, message, message)
	} // for
}

func BenchmarkRavenCParser(b *testing.B) {
	message := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	builder, _ := RavenCBuilder(1, 1, message, message, message)

	for i := 0; i < b.N; i++ {
		_, _ = RavenCParser[*msgs.RavenTest, *msgs.RavenTest](builder)
	} // for
}
