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

func (this *SuiteTCPSession) TestStartStop() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	go testerl.get().Start(SessionID(0), newCoderTester(true, true), newReactorTester(true))
	go testerc.get().Start(SessionID(1), newCoderTester(true, true), newReactorTester(true))

	time.Sleep(this.timeout)
	testerl.get().Stop()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	testerReactorl := newReactorTester(true)
	go testerl.get().Start(SessionID(0), newCoderTester(true, true), testerReactorl)
	testerReactorc := newReactorTester(true)
	go testerc.get().Start(SessionID(1), newCoderTester(true, true), testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.True(this.T(), testerReactorc.valid(this.message))

	time.Sleep(this.timeout)
	testerl.get().Send("!?")
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.valid(this.message))

	time.Sleep(this.timeout)
	testerl.get().StopWait()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	testerCoderl := newCoderTester(false, true)
	testerReactorl := newReactorTester(true)
	go testerl.get().Start(SessionID(0), testerCoderl, testerReactorl)
	testerCoderc := newCoderTester(true, true)
	testerReactorc := newReactorTester(true)
	go testerc.get().Start(SessionID(1), testerCoderc, testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorl.valid(nil))

	time.Sleep(this.timeout)
	testerl.get().StopWait()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	testerCoderl := newCoderTester(true, true)
	testerReactorl := newReactorTester(true)
	go testerl.get().Start(SessionID(0), testerCoderl, testerReactorl)
	testerCoderc := newCoderTester(true, false)
	testerReactorc := newReactorTester(true)
	go testerc.get().Start(SessionID(1), testerCoderc, testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.valid(nil))

	time.Sleep(this.timeout)
	testerl.get().StopWait()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestReceiveFailed() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	testerCoderl := newCoderTester(true, true)
	testerReactorl := newReactorTester(true)
	go testerl.get().Start(SessionID(0), testerCoderl, testerReactorl)
	testerCoderc := newCoderTester(true, true)
	testerReactorc := newReactorTester(false)
	go testerc.get().Start(SessionID(1), testerCoderc, testerReactorc)

	time.Sleep(this.timeout)
	testerl.get().Send(this.message)
	time.Sleep(this.timeout)
	assert.False(this.T(), testerReactorc.valid(nil))

	time.Sleep(this.timeout)
	testerl.get().StopWait()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestTCPSession() {
	testerl := newSessionTester()
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(testerl.complete)

	testerc := newSessionTester()
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(testerc.complete)

	assert.True(this.T(), testerl.wait())
	assert.True(this.T(), testerl.valid())
	assert.True(this.T(), testerc.wait())
	assert.True(this.T(), testerc.valid())

	go testerl.get().Start(SessionID(0), newCoderTester(true, true), newReactorTester(true))
	go testerc.get().Start(SessionID(1), newCoderTester(true, true), newReactorTester(true))

	time.Sleep(this.timeout)
	assert.Equal(this.T(), SessionID(0), testerl.get().SessionID())
	assert.NotNil(this.T(), testerl.get().RemoteAddr())
	assert.NotNil(this.T(), testerl.get().LocalAddr())
	assert.Equal(this.T(), SessionID(1), testerc.get().SessionID())
	assert.NotNil(this.T(), testerc.get().RemoteAddr())
	assert.NotNil(this.T(), testerc.get().LocalAddr())

	time.Sleep(this.timeout)
	testerl.get().StopWait()
	testerc.get().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}
