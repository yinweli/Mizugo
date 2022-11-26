package nets

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/goleak"

	"github.com/yinweli/Mizugo/mizugo/utils"
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
}

func (this *SuiteTCPSession) SetupSuite() {
	this.Change("test-tcpSession")
	this.ip = ""
	this.port = "3002"
	this.timeout = time.Second
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
	sessionl := newTestSession(this.timeout, "session server")
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession(this.timeout, "session client")
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl.Receive, sessionl.Inform)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc.Receive, sessionc.Inform)

	time.Sleep(this.timeout)
	sessionl.Session().Stop()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	sessionl := newTestSession(this.timeout, "session server")
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession(this.timeout, "session client")
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl.Receive, sessionl.Inform)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc.Receive, sessionc.Inform)

	time.Sleep(this.timeout)
	sessionl.Session().Send([]byte("send packet"))
	time.Sleep(this.timeout)
	assert.Equal(this.T(), "send packet", string(sessionc.Packet()))

	time.Sleep(this.timeout)
	sessionl.Session().StopWait()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSessionID() {
	sessionl := newTestSession(this.timeout, "session server")
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession(this.timeout, "session client")
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl.Receive, sessionl.Inform)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc.Receive, sessionc.Inform)

	time.Sleep(this.timeout)
	assert.Equal(this.T(), SessionID(0), sessionl.Session().SessionID())
	assert.Equal(this.T(), SessionID(1), sessionc.Session().SessionID())

	time.Sleep(this.timeout)
	sessionl.Session().Stop()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func newTestSession(timeout time.Duration, name string) *testSession {
	return &testSession{
		timeout: utils.NewWaitTimeout(timeout),
		name:    name,
	}
}

type testSession struct {
	lock    sync.Mutex
	timeout *utils.WaitTimeout
	name    string
	session Sessioner
	packet  []byte
	err     error
}

func (this *testSession) Wait() bool {
	return this.timeout.Wait()
}

func (this *testSession) Complete(session Sessioner, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err == nil {
		this.session = session
		fmt.Printf("%s remote addr: %s\n", this.name, session.RemoteAddr().String())
		fmt.Printf("%s local addr: %s\n", this.name, session.LocalAddr().String())
	} else {
		this.err = err
		fmt.Printf("%s: %s\n", this.name, err.Error())
	} // if

	this.timeout.Done()
}

func (this *testSession) Receive(packet []byte) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.packet = packet
}

func (this *testSession) Inform(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.err = err
}

func (this *testSession) Session() Sessioner {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session
}

func (this *testSession) Packet() []byte {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.packet
}

func (this *testSession) Error() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.err
}
