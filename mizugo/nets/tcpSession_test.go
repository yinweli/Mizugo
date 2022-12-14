package nets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/testdata"
)

func TestTCPSession(t *testing.T) {
	suite.Run(t, new(SuiteTCPSession))
}

type SuiteTCPSession struct {
	suite.Suite
	testdata.TestEnv
	ip      string
	port    string
	timeout time.Duration
	message string
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Change("test-nets-tcpSession")
	this.ip = ""
	this.port = "3002"
	this.timeout = time.Second
	this.message = "message"
}

func (this *SuiteTCPSession) TearDownSuite() {
	this.Restore()
}

func (this *SuiteTCPSession) TearDownTest() {
	goleak.VerifyNone(this.T())
}

func (this *SuiteTCPSession) TestNewTCPSession() {
	assert.NotNil(this.T(), NewTCPSession(nil))
}

func (this *SuiteTCPSession) TestStart() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	go testerl.get().Start(SessionID(0), newSessionTester(true, true, true))
	go testerc.get().Start(SessionID(1), newSessionTester(true, true, true))

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().Stop() // 這裡故意用一般結束來測試看看
	testerc.get().StopWait()
}

func (this *SuiteTCPSession) TestSend() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	testerReactorl := newSessionTester(true, true, true)
	go testerl.get().Start(SessionID(0), testerReactorl)
	testerReactorc := newSessionTester(true, true, true)
	go testerc.get().Start(SessionID(1), testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.True(this.T(), testerReactorc.validMessage(this.message))

	time.Sleep(this.timeout)
	testerl.get().Send("!?")
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.validMessage(this.message))

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	testerReactorl := newSessionTester(false, true, true)
	go testerl.get().Start(SessionID(0), testerReactorl)
	testerReactorc := newSessionTester(true, true, true)
	go testerc.get().Start(SessionID(1), testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorl.validError())

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	testerReactorl := newSessionTester(true, true, true)
	go testerl.get().Start(SessionID(0), testerReactorl)
	testerReactorc := newSessionTester(true, false, true)
	go testerc.get().Start(SessionID(1), testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.validError())

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteTCPSession) TestReceiveFailed() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	testerReactorl := newSessionTester(true, true, true)
	go testerl.get().Start(SessionID(0), testerReactorl)
	testerReactorc := newSessionTester(true, true, false)
	go testerc.get().Start(SessionID(1), testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.validError())

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()
}

func (this *SuiteTCPSession) TestTCPSession() {
	testerl := newCompleteTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(testerl)

	testerc := newCompleteTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	client.Connect(testerc)

	time.Sleep(this.timeout)
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.valid())

	go testerl.get().Start(SessionID(0), newSessionTester(true, true, true))
	go testerc.get().Start(SessionID(1), newSessionTester(true, true, true))

	time.Sleep(this.timeout)
	assert.Equal(this.T(), SessionID(0), testerl.get().SessionID())
	assert.NotNil(this.T(), testerl.get().RemoteAddr())
	assert.NotNil(this.T(), testerl.get().LocalAddr())
	assert.Equal(this.T(), SessionID(1), testerc.get().SessionID())
	assert.NotNil(this.T(), testerc.get().RemoteAddr())
	assert.NotNil(this.T(), testerc.get().LocalAddr())

	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
	testerl.get().StopWait()
	testerc.get().StopWait()
}
