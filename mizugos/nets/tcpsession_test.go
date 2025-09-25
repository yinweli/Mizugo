package nets

import (
	"testing"

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
	this.NotNil(NewTCPSession(nil))

	testl := newTestNet(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())

	trials.WaitTimeout()
	this.NotNil(testl.Get().RemoteAddr())
	this.NotNil(testl.Get().LocalAddr())
	this.NotNil(testc.Get().RemoteAddr())
	this.NotNil(testc.Get().LocalAddr())

	trials.WaitTimeout()
	owner := "owner"
	testc.Get().SetOwner(owner)
	this.Equal(owner, testc.Get().GetOwner())

	trials.WaitTimeout()
	testc.Get().StopWait()
	this.Nil(listen.Stop())
}

func (this *SuiteTCPSession) TestStart() {
	addr := host{port: "9003"}
	testl := newTestNet(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc1 := newTestNet(true, true, true)
	client1 := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client1.Connect(testc1.Bind, testc1.Unbind, testc1.Wrong)

	testc2 := newTestNet(true, true, true)
	client2 := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client2.Connect(testc2.Bind, testc2.Unbind, testc2.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testl.ValidBind())
	this.True(testc1.Valid())
	this.True(testc1.ValidBind())
	this.True(testc1.ValidSession())
	this.True(testc2.Valid())
	this.True(testc2.ValidBind())
	this.True(testc2.ValidSession())

	trials.WaitTimeout()
	testc1.Get().Stop()
	testc2.Get().StopWait()
	this.Nil(listen.Stop())

	trials.WaitTimeout()
	this.True(testc1.ValidUnbind())
	this.True(testc1.ValidStart())
	this.True(testc1.ValidStop())
	this.False(testc1.ValidSession())
	this.True(testc2.ValidUnbind())
	this.True(testc2.ValidStart())
	this.True(testc2.ValidStop())
	this.False(testc2.ValidSession())
}

func (this *SuiteTCPSession) TestStartFailed() {
	addr := host{port: "9004"}
	testl := newTestNet(false, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.ValidUnbind())

	trials.WaitTimeout()
	this.Nil(listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	addr := host{port: "9005"}
	message := "message"
	testl := newTestNet(true, true, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send(message)
	trials.WaitTimeout()
	this.True(testl.ValidEncode())
	this.True(testl.ValidSend())
	this.True(testc.ValidDecode())
	this.True(testc.ValidRecv())
	this.True(testc.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().Send(message)
	trials.WaitTimeout()
	this.True(testc.ValidEncode())
	this.True(testc.ValidSend())
	this.True(testl.ValidDecode())
	this.True(testl.ValidRecv())
	this.True(testl.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().Send("")
	trials.WaitTimeout()
	this.True(testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send("!?")
	trials.WaitTimeout()
	this.False(testc.ValidMessage(message))

	trials.WaitTimeout()
	testc.Get().StopWait()
	this.Nil(listen.Stop())
}

func (this *SuiteTCPSession) TestSendFailed() {
	addr := host{port: "9006"}
	message := "message"
	testl := newTestNet(true, false, true)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())

	trials.WaitTimeout()
	testl.Get().Send(message)
	trials.WaitTimeout()
	this.False(testl.Valid())

	trials.WaitTimeout()
	this.Nil(listen.Stop()) // 因為編碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}

func (this *SuiteTCPSession) TestRecvFailed() {
	addr := host{port: "9007"}
	message := "message"
	testl := newTestNet(true, true, false)
	listen := NewTCPListen(addr.ip, addr.port)
	listen.Listen(testl.Bind, testl.Unbind, testl.Wrong)

	testc := newTestNet(true, true, true)
	client := NewTCPConnect(addr.ip, addr.port, trials.Timeout)
	client.Connect(testc.Bind, testc.Unbind, testc.Wrong)

	trials.WaitTimeout()
	this.True(testl.Valid())
	this.True(testc.Valid())

	trials.WaitTimeout()
	testc.Get().Send(message)
	trials.WaitTimeout()
	this.False(testl.Valid())

	trials.WaitTimeout()
	this.Nil(listen.Stop()) // 因為編碼失敗, 會直接導致連接中斷, 所以不必關閉客戶端連接
}
