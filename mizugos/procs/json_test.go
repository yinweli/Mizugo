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
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	target := NewJson()
	assert.NotNil(this.T(), target)
	input, _ := JsonMarshal(messageID, message)

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
	assert.Equal(this.T(), input, decode)

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteJson) TestProcess() {
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	target := NewJson()
	input, _ := JsonMarshal(messageID, message)

	valid := false
	target.Add(messageID, func(message any) {
		valid = assert.Equal(this.T(), input, message)
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), valid)

	input, _ = JsonMarshal(0, message)
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuiteJson) TestMarshal() {
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	output1, err := JsonMarshal(messageID, message)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)

	_, err = JsonMarshal(messageID, nil)
	assert.NotNil(this.T(), err)

	messageID, output2, err := JsonUnmarshal[msgs.JsonTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), messageID, messageID)
	assert.Equal(this.T(), message, output2)

	_, _, err = JsonUnmarshal[msgs.JsonTest](nil)
	assert.NotNil(this.T(), err)

	_, _, err = JsonUnmarshal[msgs.JsonTest]("!?")
	assert.NotNil(this.T(), err)
}

func BenchmarkJsonEncode(b *testing.B) {
	target := NewJson()
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkJsonDecode(b *testing.B) {
	target := NewJson()
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkJsonMarshal(b *testing.B) {
	input := &msgs.JsonTest{
		Data: "benchmark marshal",
	}

	for i := 0; i < b.N; i++ {
		_, _ = JsonMarshal(1, input)
	} // for
}

func BenchmarkJsonUnmarshal(b *testing.B) {
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark unmarshal",
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = JsonUnmarshal[msgs.JsonTest](input)
	} // for
}
