package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/mizugo/cryptos"
	"github.com/yinweli/Mizugo/mizugo/msgs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestJson(t *testing.T) {
	suite.Run(t, new(SuiteJson))
}

type SuiteJson struct {
	suite.Suite
	testdata.Env
}

func (this *SuiteJson) SetupSuite() {
	this.Env = testdata.EnvSetup("test-procs-json")
}

func (this *SuiteJson) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteJson) TearDownTest() {
	testdata.Leak(this.T(), true)
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

func (this *SuiteJson) TestEncodeBase64() {
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	target := NewJson().Base64(true)
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

func (this *SuiteJson) TestEncodeDesCBC() {
	key := cryptos.RandDesKeyString()
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	target := NewJson().DesCBC(true, key, key) // 這裡偷懶把key跟iv都設為key
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

func (this *SuiteJson) TestEncodeAll() {
	key := cryptos.RandDesKeyString()
	messageID := MessageID(1)
	message := &msgs.JsonTest{
		Data: "json test",
	}
	target := NewJson().Base64(true).DesCBC(true, key, key) // 這裡偷懶把key跟iv都設為key
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
