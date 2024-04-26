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
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object, object)
	target := NewRaven()
	assert.NotNil(this.T(), target)
	encode, err := target.Encode(output)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)
	assert.IsType(this.T(), []byte{}, encode)
	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestDecode() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output1, _ := RavenSBuilder(1, object, object)
	output2, _ := proto.Marshal(output1.(proto.Message))
	target := NewRaven()
	decode, err := target.Decode(output2)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(output1.(proto.Message), decode.(proto.Message)))
	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = target.Decode([]byte(testdata.Unknown))
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestProcess() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenSBuilder(1, object, object)
	valid := false
	target := NewRaven()
	target.Add(1, func(message any) {
		valid = proto.Equal(output.(proto.Message), message.(proto.Message))
	})
	assert.Nil(this.T(), target.Process(output))
	assert.True(this.T(), valid)
	output, _ = RavenSBuilder(2, object, object)
	assert.NotNil(this.T(), target.Process(output))
	assert.NotNil(this.T(), target.Process(nil))
	assert.NotNil(this.T(), target.Process(testdata.Unknown))
}

func (this *SuiteRaven) TestEncodeClient() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenSBuilder(1, object, object)
	target := NewRavenClient()
	assert.NotNil(this.T(), target)
	encode, err := target.Encode(output)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)
	assert.IsType(this.T(), []byte{}, encode)
	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestDecodeClient() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output1, _ := RavenCBuilder(1, 1, object, object, object)
	output2, _ := proto.Marshal(output1.(proto.Message))
	target := NewRavenClient()
	decode, err := target.Decode(output2)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.True(this.T(), proto.Equal(output1.(proto.Message), decode.(proto.Message)))
	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = target.Decode([]byte(testdata.Unknown))
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestProcessClient() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object)
	valid := false
	target := NewRavenClient()
	target.Add(1, func(message any) {
		valid = proto.Equal(output.(proto.Message), message.(proto.Message))
	})
	assert.Nil(this.T(), target.Process(output))
	assert.True(this.T(), valid)
	output, _ = RavenCBuilder(2, 1, object, object, object)
	assert.NotNil(this.T(), target.Process(output))
	assert.NotNil(this.T(), target.Process(nil))
	assert.NotNil(this.T(), target.Process(testdata.Unknown))
}

func (this *SuiteRaven) TestRavenS() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output1, err := RavenSBuilder(1, object, object)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)
	_, err = RavenSBuilder(1, nil, object)
	assert.NotNil(this.T(), err)
	_, err = RavenSBuilder(1, object, nil)
	assert.NotNil(this.T(), err)
	output2, err := RavenSParser[msgs.RavenTest, msgs.RavenTest](output1)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output2)
	assert.NotZero(this.T(), output2.Size())
	assert.NotNil(this.T(), output2.Detail())
	assert.Equal(this.T(), int32(1), output2.MessageID)
	assert.True(this.T(), proto.Equal(object, output2.Header))
	assert.True(this.T(), proto.Equal(object, output2.Request))
	_, err = RavenSParser[msgs.RavenTest, msgs.RavenTest](nil)
	assert.NotNil(this.T(), err)
	_, err = RavenSParser[msgs.RavenTest, msgs.RavenTest](testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = RavenSParser[int, msgs.RavenTest](output1)
	assert.NotNil(this.T(), err)
	_, err = RavenSParser[msgs.RavenTest, int](output1)
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestRavenC() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output1, err := RavenCBuilder(1, 1, object, object, object, object)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)
	_, err = RavenCBuilder(1, 1, nil, object)
	assert.NotNil(this.T(), err)
	_, err = RavenCBuilder(1, 1, object, nil)
	assert.NotNil(this.T(), err)
	_, err = RavenCBuilder(1, 1, object, object, nil)
	assert.NotNil(this.T(), err)
	output2, err := RavenCParser[msgs.RavenTest, msgs.RavenTest](output1)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output2)
	assert.NotZero(this.T(), output2.Size())
	assert.NotNil(this.T(), output2.Detail())
	assert.Equal(this.T(), int32(1), output2.MessageID)
	assert.Equal(this.T(), int32(1), output2.ErrID)
	assert.True(this.T(), proto.Equal(object, output2.Header))
	assert.True(this.T(), proto.Equal(object, output2.Request))
	assert.True(this.T(), proto.Equal(object, RavenRespond[msgs.RavenTest](output2.Respond)))
	assert.Nil(this.T(), RavenRespond[msgs.RavenTest](nil))
	assert.Nil(this.T(), RavenRespond[int](output2.Respond))
	assert.True(this.T(), proto.Equal(object, RavenRespondAt[msgs.RavenTest](output2.Respond, 1)))
	assert.Nil(this.T(), RavenRespondAt[msgs.RavenTest](output2.Respond, 10))
	assert.Nil(this.T(), RavenRespondAt[msgs.RavenTest](nil, 1))
	assert.Nil(this.T(), RavenRespondAt[int](output2.Respond, 1))
	_, err = RavenCParser[msgs.RavenTest, msgs.RavenTest](nil)
	assert.NotNil(this.T(), err)
	_, err = RavenCParser[msgs.RavenTest, msgs.RavenTest](testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = RavenCParser[int, msgs.RavenTest](output1)
	assert.NotNil(this.T(), err)
	_, err = RavenCParser[msgs.RavenTest, int](output1)
	assert.NotNil(this.T(), err)
}

func (this *SuiteRaven) TestRavenTest() {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object, object)
	assert.True(this.T(), RavenTestMessageID(output, 1))
	assert.False(this.T(), RavenTestMessageID(output, 0))
	assert.False(this.T(), RavenTestMessageID(nil, 1))
	assert.False(this.T(), RavenTestMessageID(testdata.Unknown, 1))
	assert.True(this.T(), RavenTestErrID(output, 1))
	assert.False(this.T(), RavenTestErrID(output, 0))
	assert.False(this.T(), RavenTestErrID(nil, 1))
	assert.False(this.T(), RavenTestErrID(testdata.Unknown, 1))
	assert.True(this.T(), RavenTestHeader(output, object))
	assert.False(this.T(), RavenTestHeader(output, &msgs.RavenTest{}))
	assert.False(this.T(), RavenTestHeader(nil, object))
	assert.False(this.T(), RavenTestHeader(testdata.Unknown, object))
	assert.True(this.T(), RavenTestRequest(output, object))
	assert.False(this.T(), RavenTestRequest(output, &msgs.RavenTest{}))
	assert.False(this.T(), RavenTestRequest(nil, object))
	assert.False(this.T(), RavenTestRequest(testdata.Unknown, object))
	assert.True(this.T(), RavenTestRespond(output, object, object))
	assert.False(this.T(), RavenTestRespond(nil, object, object))
	assert.False(this.T(), RavenTestRespond(testdata.Unknown, object, object))
	assert.False(this.T(), RavenTestRespond(output, object, object, object))
	assert.False(this.T(), RavenTestRespond(output, &msgs.RavenTest{}))
	assert.True(this.T(), RavenTestRespondType(output, object, object))
	assert.False(this.T(), RavenTestRespondType(nil, object, object))
	assert.False(this.T(), RavenTestRespondType(testdata.Unknown, object, object))
	assert.False(this.T(), RavenTestRespondType(output, object, object, object))
	assert.False(this.T(), RavenTestRespondType(output, &msgs.ProtoTest{}))
	assert.True(this.T(), RavenTestRespondLength(output, 2))
	assert.False(this.T(), RavenTestRespondLength(output, 3))
	assert.False(this.T(), RavenTestRespondLength(nil, 2))
	assert.False(this.T(), RavenTestRespondLength(testdata.Unknown, 2))
}

func BenchmarkRavenEncode(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object)
	target := NewRaven()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(output)
	} // for
}

func BenchmarkRavenDecode(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenSBuilder(1, object, object)
	target := NewRaven()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(output)
	} // for
}

func BenchmarkRavenEncodeClient(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenSBuilder(1, object, object)
	target := NewRavenClient()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(output)
	} // for
}

func BenchmarkRavenDecodeClient(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object)
	target := NewRavenClient()

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(output)
	} // for
}

func BenchmarkRavenSBuilder(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = RavenSBuilder(1, object, object)
	} // for
}

func BenchmarkRavenSParser(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenSBuilder(1, object, object)

	for i := 0; i < b.N; i++ {
		_, _ = RavenSParser[msgs.RavenTest, msgs.RavenTest](output)
	} // for
}

func BenchmarkRavenCBuilder(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = RavenCBuilder(1, 1, object, object, object)
	} // for
}

func BenchmarkRavenCParser(b *testing.B) {
	object := &msgs.RavenTest{
		Data: testdata.Unknown,
	}
	output, _ := RavenCBuilder(1, 1, object, object, object)

	for i := 0; i < b.N; i++ {
		_, _ = RavenCParser[msgs.RavenTest, msgs.RavenTest](output)
	} // for
}
