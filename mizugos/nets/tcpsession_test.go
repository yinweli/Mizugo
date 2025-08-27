package nets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
	"github.com/yinweli/Mizugo/v2/testdata"
)

func TestTCPSession(t *testing.T) {
	suite.Run(t, new(SuiteTCPSession))
}

type SuiteTCPSession struct {
	suite.Suite
	trials.Catalog
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Catalog = trials.Prepare(testdata.PathWork("test-nets-tcpSession"))
}

func (this *SuiteTCPSession) TearDownSuite() {
	trials.Restore(this.Catalog)
}

func (this *SuiteTCPSession) TestTCPSession() {
	addr := host{port: "9002"}
	assert.NotNil(this.T(), NewTCPSession(nil))

	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testc.Valid())

	trials.WaitTimeout()
	assert.NotNil(this.T(), testl.Get().RemoteAddr())
	assert.NotNil(this.T(), testl.Get().LocalAddr())
	assert.NotNil(this.T(), testc.Get().RemoteAddr())
	assert.NotNil(this.T(), testc.Get().LocalAddr())

	trials.WaitTimeout()
	owner := "owner"
	testc.Get().SetOwner(owner)
	assert.Equal(this.T(), owner, testc.Get().GetOwner())

	trials.WaitTimeout()
	testc.Get().StopWait()
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestStart() {
	addr := host{port: "9003"}
	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc1 := newTester(true, true, true)
	client1 := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client1.Connect(testc1.Bind, testc1.Unbind, testc1.Wrong)

	testc2 := newTester(true, true, true)
	client2 := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client2.Connect(testc2.Bind, testc2.Unbind, testc2.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testl.ValidBind())
	assert.True(this.T(), testc1.Valid())
	assert.True(this.T(), testc1.ValidBind())
	assert.True(this.T(), testc1.ValidSession())
	assert.True(this.T(), testc2.Valid())
	assert.True(this.T(), testc2.ValidBind())
	assert.True(this.T(), testc2.ValidSession())

	trials.WaitTimeout()
	testc1.Get().Stop()
	testc2.Get().StopWait()
	assert.Nil(this.T(), listen.Stop())

	trials.WaitTimeout()
	assert.True(this.T(), testc1.ValidUnbind())
	assert.True(this.T(), testc1.ValidStart())
	assert.True(this.T(), testc1.ValidStop())
	assert.False(this.T(), testc1.ValidSession())
	assert.True(this.T(), testc2.ValidUnbind())
	assert.True(this.T(), testc2.ValidStart())
	assert.True(this.T(), testc2.ValidStop())
	assert.False(this.T(), testc2.ValidSession())
}

func (this *SuiteTCPSession) TestStartFailed() {
	addr := host{port: "9004"}
	testl := newTester(false, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.ValidUnbind())

	trials.WaitTimeout()
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	addr := host{port: "9005"}
	message := "message"
	testl := newTester(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send(message)
	trials.WaitTimeout()
	assert.True(this.T(), testl.ValidEncode())
	assert.True(this.T(), testl.ValidSend())
	assert.True(this.T(), testc.ValidDecode())
	assert.True(this.T(), testc.ValidRecv())
	assert.True(this.T(), testc.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().Send(message)
	trials.WaitTimeout()
	assert.True(this.T(), testc.ValidEncode())
	assert.True(this.T(), testc.ValidSend())
	assert.True(this.T(), testl.ValidDecode())
	assert.True(this.T(), testl.ValidRecv())
	assert.True(this.T(), testl.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().Send("")
	trials.WaitTimeout()
	assert.True(this.T(), testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send("!?")
	trials.WaitTimeout()
	assert.False(this.T(), testc.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().StopWait()
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSendFailed() {
	addr := host{port: "9006"}
	message := "message"
	testl := newTester(true, false, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send(message)
	trials.WaitTimeout()
	assert.False(this.T(), testl.Valid())

	trials.WaitTimeout()
	assert.Nil(this.T(), listen.Stop()) // 因為編碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}

func (this *SuiteTCPSession) TestRecvFailed() {
	addr := host{port: "9007"}
	message := "message"
	testl := newTester(true, true, false)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTester(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	assert.True(this.T(), testl.Valid())
	assert.True(this.T(), testc.Valid())

	trials.WaitTimeout()
	testc.Get().Send(message)
	trials.WaitTimeout()
	assert.False(this.T(), testl.Valid())

	trials.WaitTimeout()
	assert.Nil(this.T(), listen.Stop()) // 因為編碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}
