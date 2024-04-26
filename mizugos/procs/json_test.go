package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/mizugos/trials"
	"github.com/yinweli/Mizugo/testdata"
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
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	assert.NotNil(this.T(), target)
	encode, err := target.Encode(output)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)
	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Encode(testdata.Unknown)
	assert.NotNil(this.T(), err)
}

func (this *SuiteJson) TestDecode() {
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	encode, _ := target.Encode(output)
	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.Equal(this.T(), output, decode)
	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)
	_, err = target.Decode(testdata.Unknown)
	assert.NotNil(this.T(), err)
	_, err = target.Decode([]byte(testdata.Unknown))
	assert.NotNil(this.T(), err)
}

func (this *SuiteJson) TestProcess() {
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	valid := false
	target := NewJson()
	target.Add(1, func(message any) {
		valid = assert.Equal(this.T(), output, message)
	})
	assert.Nil(this.T(), target.Process(output))
	assert.True(this.T(), valid)
	output, _ = JsonMarshal(0, &msgs.JsonTest{})
	assert.NotNil(this.T(), target.Process(output))
	assert.NotNil(this.T(), target.Process(nil))
	assert.NotNil(this.T(), target.Process(testdata.Unknown))
}

func (this *SuiteJson) TestMarshal() {
	output, err := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output)
	_, err = JsonMarshal(1, nil)
	assert.NotNil(this.T(), err)
}

func (this *SuiteJson) TestUnmarshal() {
	object := &msgs.JsonTest{
		Data: testdata.Unknown,
	}
	output1, _ := JsonMarshal(1, object)
	messageID, output2, err := JsonUnmarshal[msgs.JsonTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), int32(1), messageID)
	assert.Equal(this.T(), object, output2)
	_, _, err = JsonUnmarshal[msgs.JsonTest](nil)
	assert.NotNil(this.T(), err)
	_, _, err = JsonUnmarshal[int](output1)
	assert.NotNil(this.T(), err)
}

func BenchmarkJsonEncode(b *testing.B) {
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(output)
	} // for
}

func BenchmarkJsonDecode(b *testing.B) {
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})
	target := NewJson()
	encode, _ := target.Encode(output)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkJsonMarshal(b *testing.B) {
	object := &msgs.JsonTest{
		Data: testdata.Unknown,
	}

	for i := 0; i < b.N; i++ {
		_, _ = JsonMarshal(1, object)
	} // for
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	output, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: testdata.Unknown,
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = JsonUnmarshal[msgs.JsonTest](output)
	} // for
}
