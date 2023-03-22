package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestJson(t *testing.T) {
	suite.Run(t, new(SuiteJson))
}

type SuiteJson struct {
	suite.Suite
	testdata.TestEnv
	key       string
	messageID MessageID
	message   *msgs.JsonTest
}

func (this *SuiteJson) SetupSuite() {
	this.TBegin("test-procs-json", "")
	this.key = cryptos.RandDesKeyString()
	this.messageID = MessageID(1)
	this.message = &msgs.JsonTest{
		Data: "json test",
	}
}

func (this *SuiteJson) TearDownSuite() {
	this.TFinal()
}

func (this *SuiteJson) TearDownTest() {
	this.TLeak(this.T(), true)
}

func (this *SuiteJson) TestNewJson() {
	assert.NotNil(this.T(), NewJson())
}

func (this *SuiteJson) TestEncode() {
	target := NewJson()
	input, _ := JsonMarshal(this.messageID, this.message)

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

func (this *SuiteJson) TestEncodeBase64() {
	target := NewJson().Base64(true)
	input, _ := JsonMarshal(this.messageID, this.message)

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

func (this *SuiteJson) TestEncodeDesCBC() {
	target := NewJson().DesCBC(true, this.key, this.key) // 這裡偷懶把key跟iv都設為key
	input, _ := JsonMarshal(this.messageID, this.message)

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

func (this *SuiteJson) TestEncodeAll() {
	target := NewJson().Base64(true).DesCBC(true, this.key, this.key) // 這裡偷懶把key跟iv都設為key
	input, _ := JsonMarshal(this.messageID, this.message)

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
	target := NewJson()
	input, _ := JsonMarshal(this.messageID, this.message)

	valid := false
	target.Add(this.messageID, func(message any) {
		valid = assert.Equal(this.T(), input, message)
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), valid)

	input, _ = JsonMarshal(0, this.message)
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuiteJson) TestMarshal() {
	output1, err := JsonMarshal(this.messageID, this.message)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)

	_, err = JsonMarshal(this.messageID, nil)
	assert.NotNil(this.T(), err)

	messageID, output2, err := JsonUnmarshal[msgs.JsonTest](output1)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), this.messageID, messageID)
	assert.Equal(this.T(), this.message, output2)

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

func BenchmarkJsonEncodeBase64(b *testing.B) {
	target := NewJson().Base64(true)
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkJsonEncodeDesCBC(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewJson().DesCBC(true, key, key)
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkJsonEncodeAll(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewJson().Base64(true).DesCBC(true, key, key)
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

func BenchmarkJsonDecodeBase64(b *testing.B) {
	target := NewJson().Base64(true)
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkJsonDecodeDesCBC(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewJson().DesCBC(true, key, key)
	input, _ := JsonMarshal(1, &msgs.JsonTest{
		Data: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkJsonDecodeAll(b *testing.B) {
	key := cryptos.RandDesKeyString()
	target := NewJson().Base64(true).DesCBC(true, key, key)
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
