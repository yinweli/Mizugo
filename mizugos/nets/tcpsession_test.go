package nets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPSession(t *testing.T) {
	suite.Run(t, new(SuiteTCPSession))
}

type SuiteTCPSession struct {
	suite.Suite
	testdata.TestEnv
	testdata.TestLeak
	host    host
	message string
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Change("test-nets-tcpSession")
	this.host = host{port: "3001"}
	this.message = "message"
}

func (this *SuiteTCPSession) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPSession) TearDownTest() {
	this.GoLeak(this.T(), true)
}

func (this *SuiteTCPSession) TestNewTCPSession() {
	assert.NotNil(this.T(), NewTCPSession(nil))
}

func (this *SuiteTCPSession) TestStart() {
	testl := newTester(true, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc1 := newTester(true, true, true, true)
	client1 := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client1.Connect(testc1.bind, testc1.unbind, testc1.wrong)

	testc2 := newTester(true, true, true, true)
	client2 := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client2.Connect(testc2.bind, testc2.unbind, testc2.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testl.validBind())
	assert.True(this.T(), testc1.valid())
	assert.True(this.T(), testc1.validBind())
	assert.True(this.T(), testc1.validSession())
	assert.True(this.T(), testc2.valid())
	assert.True(this.T(), testc2.validBind())
	assert.True(this.T(), testc2.validSession())

	time.Sleep(testdata.Timeout)
	testc1.get().Stop()
	testc2.get().StopWait()
	assert.Nil(this.T(), listen.Stop())

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testc1.validUnbind())
	assert.False(this.T(), testc1.validSession())
	assert.True(this.T(), testc2.validUnbind())
	assert.False(this.T(), testc2.validSession())
}

func (this *SuiteTCPSession) TestStartFailed() {
	testl := newTester(false, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true, true)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.False(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	testl := newTester(true, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true, true)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.validEncode())
	assert.True(this.T(), testl.validAfterSend())
	assert.True(this.T(), testc.validDecode())
	assert.True(this.T(), testc.validReceive())
	assert.True(this.T(), testc.validAfterRecv())
	assert.True(this.T(), testc.validMessage(this.message))

	time.Sleep(testdata.Timeout)
	testc.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testc.validEncode())
	assert.True(this.T(), testc.validAfterSend())
	assert.True(this.T(), testl.validDecode())
	assert.True(this.T(), testl.validReceive())
	assert.True(this.T(), testl.validAfterRecv())
	assert.True(this.T(), testl.validMessage(this.message))

	time.Sleep(testdata.Timeout)
	testc.get().Send("")
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send("!?")
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testc.validMessage(this.message))

	time.Sleep(testdata.Timeout)
	testc.get().StopWait()
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	testl := newTester(true, false, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true, true)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop()) // ??????????????????, ???????????????????????????, ?????????????????????????????????
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	testl := newTester(true, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, false, true)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop()) // ??????????????????, ???????????????????????????, ?????????????????????????????????
}

func (this *SuiteTCPSession) TestReceiveFailed() {
	testl := newTester(true, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true, false)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop()) // ??????????????????, ???????????????????????????, ?????????????????????????????????
}

func (this *SuiteTCPSession) TestTCPSession() {
	testl := newTester(true, true, true, true)
	listen := NewTCPListen(this.host.ip, this.host.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true, false)
	client := NewTCPConnect(this.host.ip, this.host.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.NotNil(this.T(), testl.get().RemoteAddr())
	assert.NotNil(this.T(), testl.get().LocalAddr())
	assert.NotNil(this.T(), testc.get().RemoteAddr())
	assert.NotNil(this.T(), testc.get().LocalAddr())

	time.Sleep(testdata.Timeout)
	owner := "owner"
	testc.get().SetOwner(owner)
	assert.Equal(this.T(), owner, testc.get().GetOwner())

	time.Sleep(testdata.Timeout)
	testc.get().StopWait()
	assert.Nil(this.T(), listen.Stop())
}
