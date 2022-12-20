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
	message string
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Change("test-nets-tcpSession")
	this.ip = ""
	this.port = "3002"
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
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(true, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, true, true)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().Stop() // 這裡故意用一般結束來測試看看
	donec.get().StopWait()
}

func (this *SuiteTCPSession) TestSend() {
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(true, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, true, true)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	donel.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.True(this.T(), bindc.validMessage(this.message))

	time.Sleep(testdata.Timeout)
	donel.get().Send("!?")
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bindc.validMessage(this.message))

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().StopWait()
	donec.get().StopWait()
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(false, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, true, true)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	donel.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bindl.validError())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().StopWait()
	donec.get().StopWait()
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(true, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, false, true)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	donel.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bindc.validError())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().StopWait()
	donec.get().StopWait()
}

func (this *SuiteTCPSession) TestReceiveFailed() {
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(true, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, true, false)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	donel.get().Send(this.message)
	time.Sleep(testdata.Timeout)
	assert.False(this.T(), bindc.validError())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().StopWait()
	donec.get().StopWait()
}

func (this *SuiteTCPSession) TestTCPSession() {
	donel := newDoneTester()
	listen := NewTCPListen(this.ip, this.port)
	listen.Listen(donel.done)

	donec := newDoneTester()
	client := NewTCPConnect(this.ip, this.port, testdata.Timeout)
	client.Connect(donec.done)

	time.Sleep(testdata.Timeout)
	assert.True(this.T(), donel.valid())
	assert.True(this.T(), donec.valid())

	bindl := newBindTester(true, true, true)
	go donel.get().Start(SessionID(0), bindl.bind)
	bindc := newBindTester(true, true, true)
	go donec.get().Start(SessionID(1), bindc.bind)

	time.Sleep(testdata.Timeout)
	assert.Equal(this.T(), SessionID(0), donel.get().SessionID())
	assert.NotNil(this.T(), donel.get().RemoteAddr())
	assert.NotNil(this.T(), donel.get().LocalAddr())
	assert.Equal(this.T(), SessionID(1), donec.get().SessionID())
	assert.NotNil(this.T(), donec.get().RemoteAddr())
	assert.NotNil(this.T(), donec.get().LocalAddr())

	time.Sleep(testdata.Timeout)
	assert.Nil(this.T(), listen.Stop())
	donel.get().StopWait()
	donec.get().StopWait()
}
