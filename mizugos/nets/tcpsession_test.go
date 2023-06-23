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
	testdata.Env
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Env = testdata.EnvSetup("test-nets-tcpSession")
}

func (this *SuiteTCPSession) TearDownSuite() {
	testdata.EnvRestore(this.Env)
}

func (this *SuiteTCPSession) TearDownTest() {
	testdata.Leak(this.T(), true)
}

func (this *SuiteTCPSession) TestTCPSession() {
	addr := host{port: "9002"}
	assert.NotNil(this.T(), NewTCPSession(nil))

	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
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

func (this *SuiteTCPSession) TestStart() {
	addr := host{port: "9002"}
	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc1 := newTester(true, true, true)
	client1 := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client1.Connect(testc1.bind, testc1.unbind, testc1.wrong)

	testc2 := newTester(true, true, true)
	client2 := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
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
	assert.True(this.T(), testc1.validStart())
	assert.True(this.T(), testc1.validStop())
	assert.False(this.T(), testc1.validSession())
	assert.True(this.T(), testc2.validUnbind())
	assert.True(this.T(), testc2.validStart())
	assert.True(this.T(), testc2.validStop())
	assert.False(this.T(), testc2.validSession())
}

func (this *SuiteTCPSession) TestStartFailed() {
	addr := host{port: "9002"}
	testl := newTester(false, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.False(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	addr := host{port: "9002"}
	message := "message"
	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(message)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.validEncode())
	assert.True(this.T(), testl.validSend())
	assert.True(this.T(), testc.validDecode())
	assert.True(this.T(), testc.validRecv())
	assert.True(this.T(), testc.validMessage(message))

	time.Sleep(testdata.Timeout)
	testc.get().Send(message)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testc.validEncode())
	assert.True(this.T(), testc.validSend())
	assert.True(this.T(), testl.validDecode())
	assert.True(this.T(), testl.validRecv())
	assert.True(this.T(), testl.validMessage(message))

	time.Sleep(testdata.Timeout)
	testc.get().Send("")
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send("!?")
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testc.validMessage(message))

	time.Sleep(testdata.Timeout)
	testc.get().StopWait()
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	addr := host{port: "9002"}
	message := "message"
	testl := newTester(true, false, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testl.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop()) // 因為編碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	addr := host{port: "9002"}
	message := "message"
	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.bind, testl.unbind, testl.wrong)

	testc := newTester(true, true, false)
	client := NewTCPConnect(addr.ip, addr.port, testdata.Timeout)
	client.Connect(testc.bind, testc.unbind, testc.wrong)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), testl.valid())
	assert.True(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	testl.get().Send(message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), testc.valid())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop()) // 因為解碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}
