package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestSimple(t *testing.T) {
	suite.Run(t, new(SuiteSimple))
}

type SuiteSimple struct {
	suite.Suite
	testdata.TestEnv
}

func (this *SuiteSimple) SetupSuite() {
	this.Change("test-procs-simple")
}

func (this *SuiteSimple) TearDownSuite() {
	this.Restore()
}

func (this *SuiteSimple) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteSimple) TestNewSimple() {
	assert.NotNil(this.T(), NewSimple())
}

func (this *SuiteSimple) TestEncodeDecode() {
	target := NewSimple()
	input, err := SimpleMarshal(1, &SimpleMsgTest{
		Message: "test encode/decode",
	})

	assert.Nil(this.T(), err)

	encode, err := target.Encode(input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), encode)

	decode, err := target.Decode(encode)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), decode)
	assert.Equal(this.T(), input, decode)

	_, err = target.Encode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode/decode"))
	assert.NotNil(this.T(), err)
}

func (this *SuiteSimple) TestProcess() {
	target := NewSimple()
	input, err := SimpleMarshal(1, &SimpleMsgTest{
		Message: "test process",
	})
	assert.Nil(this.T(), err)

	valid := false
	target.Add(input.MessageID, func(message any) {
		valid = assert.Equal(this.T(), input, message)
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), valid)

	input.MessageID = 2
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuiteSimple) TestMarshal() {
	inputID := MessageID(1)
	input := &SimpleMsgTest{
		Message: "test marshal/unmarshal",
	}

	message, err := SimpleMarshal(inputID, input)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), message)

	outputID, output, err := SimpleUnmarshal[SimpleMsgTest](message)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), inputID, outputID)
	assert.Equal(this.T(), input, output)

	_, _, err = SimpleUnmarshal[SimpleMsgTest](nil)
	assert.NotNil(this.T(), err)

	_, _, err = SimpleUnmarshal[int64](message)
	assert.NotNil(this.T(), err)
}

func BenchmarkSimpleEncode(b *testing.B) {
	target := NewSimple()
	input, _ := SimpleMarshal(1, &SimpleMsgTest{
		Message: "benchmark encode",
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkSimpleDecode(b *testing.B) {
	target := NewSimple()
	input, _ := SimpleMarshal(1, &SimpleMsgTest{
		Message: "benchmark decode",
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkSimpleMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = SimpleMarshal(1, &SimpleMsgTest{
			Message: "benchmark marshal",
		})
	} // for
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	input, _ := SimpleMarshal(1, &SimpleMsgTest{
		Message: "benchmark marshal",
	})

	for i := 0; i < b.N; i++ {
		_, _, _ = SimpleUnmarshal[SimpleMsgTest](input)
	} // for
}
