package procs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"

	"github.com/yinweli/Mizugo/mizugos/cryptos"
	"github.com/yinweli/Mizugo/mizugos/msgs"
	"github.com/yinweli/Mizugo/testdata"
)

func TestPList(t *testing.T) {
	suite.Run(t, new(SuitePList))
}

type SuitePList struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	key       []byte
	messageID MessageID
	message   *msgs.PListTest
}

func (this *SuitePList) SetupSuite() {
	this.Change("test-procs-plist")
	this.key = cryptos.RandDesKey()
	this.messageID = MessageID(1)
	this.message = &msgs.PListTest{
		Data: "plist test",
	}
}

func (this *SuitePList) TearDownSuite() {
	this.Restore()
}

func (this *SuitePList) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuitePList) TestNewPList() {
	assert.NotNil(this.T(), NewPList())
}

func (this *SuitePList) TestEncode() {
	target := NewPList().Key(this.key)
	input := MarshalPListMsg([]TestMsg{
		{MessageID: this.messageID, Message: this.message},
	})

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
	assert.True(this.T(), proto.Equal(input, decode.(*msgs.PListMsg)))

	_, err = target.Decode(nil)
	assert.NotNil(this.T(), err)

	_, err = target.Decode([]byte("unknown encode"))
	assert.NotNil(this.T(), err)
}

func (this *SuitePList) TestProcess() {
	target := NewPList().Key(this.key)
	input := MarshalPListMsg([]TestMsg{
		{MessageID: this.messageID, Message: this.message},
	})

	validSend := false
	target.Send(func(message any) {
		_, validSend = message.(*msgs.PListMsg)
	})
	validProcess := false
	target.Add(this.messageID, func(message any) {
		if context, ok := message.(*PListContext); ok {
			messageID, msg, err := PListUnmarshal[msgs.PListTest](context)
			validProcess = err == nil && this.messageID == messageID && proto.Equal(this.message, msg)
		} // if
	})
	assert.Nil(this.T(), target.Process(input))
	assert.True(this.T(), validSend)
	assert.True(this.T(), validProcess)

	input = MarshalPListMsg([]TestMsg{
		{MessageID: 0, Message: this.message},
	})
	assert.NotNil(this.T(), target.Process(input))

	assert.NotNil(this.T(), target.Process(nil))
}

func (this *SuitePList) TestPListContext() {
	target := &PListContext{
		request: MarshalPListMsg([]TestMsg{
			{MessageID: 1, Message: this.message},
			{MessageID: 2, Message: this.message},
			{MessageID: 3, Message: this.message},
		}).Messages,
	}
	assert.True(this.T(), target.next())
	assert.Equal(this.T(), MessageID(1), target.messageID())
	assert.NotNil(this.T(), target.message())
	assert.True(this.T(), target.next())
	assert.Equal(this.T(), MessageID(2), target.messageID())
	assert.NotNil(this.T(), target.message())
	assert.True(this.T(), target.next())
	assert.Equal(this.T(), MessageID(3), target.messageID())
	assert.NotNil(this.T(), target.message())
	assert.False(this.T(), target.next())

	target = &PListContext{}
	assert.Nil(this.T(), target.AddRespond(1, this.message))
	assert.Nil(this.T(), target.AddRespond(2, this.message))
	assert.Nil(this.T(), target.AddRespond(3, this.message))
}

func (this *SuitePList) TestMarshal() {
	target := &PListContext{}
	assert.Nil(this.T(), target.AddRespond(1, this.message))
	assert.Nil(this.T(), target.AddRespond(2, this.message))
	assert.Nil(this.T(), target.AddRespond(3, this.message))
	output1, err := PListMarshal(target)
	assert.Nil(this.T(), err)
	assert.NotNil(this.T(), output1)

	_, err = PListMarshal(nil)
	assert.NotNil(this.T(), err)

	target = &PListContext{
		request: MarshalPListMsg([]TestMsg{
			{MessageID: 1, Message: this.message},
			{MessageID: 2, Message: this.message},
			{MessageID: 3, Message: this.message},
		}).Messages,
	}
	target.next()
	messageID, output2, err := PListUnmarshal[msgs.PListTest](target)
	assert.Nil(this.T(), err)
	assert.Equal(this.T(), MessageID(1), messageID)
	assert.True(this.T(), proto.Equal(this.message, output2))

	_, _, err = PListUnmarshal[msgs.PListTest](nil)
	assert.NotNil(this.T(), err)
}

func BenchmarkPListEncode(b *testing.B) {
	target := NewPList()
	input := MarshalPListMsg([]TestMsg{
		{MessageID: 1, Message: &msgs.PListTest{Data: "benchmark encode"}},
	})

	for i := 0; i < b.N; i++ {
		_, _ = target.Encode(input)
	} // for
}

func BenchmarkPListDecode(b *testing.B) {
	target := NewPList()
	input := MarshalPListMsg([]TestMsg{
		{MessageID: 1, Message: &msgs.PListTest{Data: "benchmark decode"}},
	})
	encode, _ := target.Encode(input)

	for i := 0; i < b.N; i++ {
		_, _ = target.Decode(encode)
	} // for
}

func BenchmarkPListMarshal(b *testing.B) {
	input := &PListContext{}
	_ = input.AddRespond(1, &msgs.PListTest{
		Data: "benchmark marshal",
	})

	for i := 0; i < b.N; i++ {
		_, _ = PListMarshal(input)
	} // for
}

func BenchmarkPListUnmarshal(b *testing.B) {
	input := &PListContext{
		request: MarshalPListMsg([]TestMsg{
			{MessageID: 1, Message: &msgs.PListTest{Data: "benchmark unmarshal"}},
		}).Messages,
	}
	input.next()

	for i := 0; i < b.N; i++ {
		_, _, _ = PListUnmarshal[msgs.PListTest](input)
	} // for
}
