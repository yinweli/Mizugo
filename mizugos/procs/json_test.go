package procs

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/msgs"
	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestJson(t *testing.T) {
	suite.Run(t, new(SuiteJson))
}

type SuiteJson struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteJson) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-procs-json"))
}

func (this *SuiteJson) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteJson) TestEncode() {
	message, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	this.NotNil(target)
	encode, err := target.Encode(message)
	this.Nil(err)
	this.NotNil(encode)
	_, err = target.Encode(nil)
	this.NotNil(err)
	_, err = target.Encode(testdata.Unknown)
	this.NotNil(err)
}

func (this *SuiteJson) TestDecode() {
	marshal, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	encode, _ := target.Encode(marshal)
	decode, err := target.Decode(encode)
	this.Nil(err)
	this.NotNil(decode)
	this.Equal(marshal, decode)
	_, err = target.Decode(nil)
	this.NotNil(err)
	_, err = target.Decode(testdata.Unknown)
	this.NotNil(err)
	_, err = target.Decode([]byte(testdata.Unknown))
	this.NotNil(err)
}

func (this *SuiteJson) TestProcess() {
	marshal, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	valid := false
	target := NewJson()
	target.Add(1, func(message any) {
		valid = this.Equal(message, message)
	})
	this.Nil(target.Process(marshal))
	this.True(valid)
	marshal, _ = JsonMarshal(0, &msgs.JsonTest{})
	this.NotNil(target.Process(marshal))
	this.NotNil(target.Process(nil))
	this.NotNil(target.Process(testdata.Unknown))
}

func (this *SuiteJson) TestMarshal() {
	marshal, err := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	this.Nil(err)
	this.NotNil(marshal)
	_, err = JsonMarshal(1, nil)
	this.NotNil(err)
}

func (this *SuiteJson) TestUnmarshal() {
	message := &msgs.JsonTest{
		Data: testdata.Unknown,
	}
	marshal, _ := JsonMarshal(1, message)
	messageID, payload, err := JsonUnmarshal[msgs.JsonTest](marshal)
	this.Nil(err)
	this.Equal(int32(1), messageID)
	this.Equal(message, payload)
	_, _, err = JsonUnmarshal[msgs.JsonTest](nil)
	this.NotNil(err)
	_, _, err = JsonUnmarshal[int](marshal)
	this.NotNil(err)
}

func BenchmarkJsonEncode(b *testing.B) {
	marshal, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(marshal)
	} // for
}

func BenchmarkJsonDecode(b *testing.B) {
	marshal, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	encode, _ := target.Encode(marshal)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkJsonMarshal(b *testing.B) {
	message := &msgs.JsonTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = JsonMarshal(1, message)
	} // for
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	marshal, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = JsonUnmarshal[msgs.JsonTest](marshal)
	} // for
}
