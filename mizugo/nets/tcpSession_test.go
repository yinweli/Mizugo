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
	this.Change("test-nets-tcpSession")
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
	sessionl := newTestSession("tcp session l - start/stop", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - start/stop", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	sessionl.Session().Stop()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSend() {
	sessionl := newTestSession("tcp session l - send", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - send", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	sessionl.Session().Send(sessionc.Name())
	time.Sleep(this.timeout)
	assert.True(this.T(), sessionc.Valid())

	time.Sleep(this.timeout)
	sessionl.Session().StopWait()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSendFailed() {
	sessionl := newTestSession("tcp session l - send failed", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - send failed", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	sessionl.Session().Send("!?")
	time.Sleep(this.timeout)
	assert.False(this.T(), sessionc.Valid())

	time.Sleep(this.timeout)
	sessionl.Session().StopWait()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestEncodeFailed() {
	sessionl := newTestSession("tcp session l - encode failed", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - encode failed", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	sessionl.EncodeFailed(true)
	sessionl.Session().Send(sessionc.Name())
	time.Sleep(this.timeout)
	assert.False(this.T(), sessionl.Valid())

	time.Sleep(this.timeout)
	sessionl.Session().StopWait()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestDecodeFailed() {
	sessionl := newTestSession("tcp session l - decode failed", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - decode failed", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	sessionc.DecodeFailed(true)
	sessionl.Session().Send(sessionc.Name())
	time.Sleep(this.timeout)
	assert.False(this.T(), sessionc.Valid())

	time.Sleep(this.timeout)
	sessionl.Session().StopWait()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func (this *SuiteTCPSession) TestSessionID() {
	sessionl := newTestSession("tcp session l - sessionID", this.timeout)
	listen := NewTCPListen(this.ip, this.port)
	go listen.Start(sessionl.Complete)

	time.Sleep(this.timeout)

	sessionc := newTestSession("tcp session c - sessionID", this.timeout)
	client := NewTCPConnect(this.ip, this.port, this.timeout)
	go client.Start(sessionc.Complete)

	assert.True(this.T(), sessionl.Wait())
	assert.NotNil(this.T(), sessionl.Session())
	go sessionl.Session().Start(SessionID(0), sessionl, sessionl)

	assert.True(this.T(), sessionc.Wait())
	assert.NotNil(this.T(), sessionc.Session())
	go sessionc.Session().Start(SessionID(1), sessionc, sessionc)

	time.Sleep(this.timeout)
	assert.Equal(this.T(), SessionID(0), sessionl.Session().SessionID())
	assert.Equal(this.T(), SessionID(1), sessionc.Session().SessionID())

	time.Sleep(this.timeout)
	sessionl.Session().Stop()
	sessionc.Session().StopWait()
	time.Sleep(this.timeout)
	assert.Nil(this.T(), listen.Stop())
}

func newTestSession(name string, timeout time.Duration) *testSession {
	return &testSession{
		name:    name,
		timeout: utils.NewWaitTimeout(timeout),
	}
}

type testSession struct {
	name         string
	timeout      *utils.WaitTimeout
	encodeFailed bool
	decodeFailed bool
	session      Sessioner
	result       error
	lock         sync.Mutex
}

func (this *testSession) Wait() bool {
	return this.timeout.Wait()
}

func (this *testSession) EncodeFailed(encodeFailed bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.encodeFailed = encodeFailed
}

func (this *testSession) DecodeFailed(decodeFailed bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.decodeFailed = decodeFailed
}

func (this *testSession) Name() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.name
}

func (this *testSession) Session() Sessioner {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session
}

func (this *testSession) Valid() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.result == nil
}

func (this *testSession) Complete(session Sessioner, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err == nil {
		this.session = session
		this.result = nil
		fmt.Printf("%s remote addr: %s\n", this.name, session.RemoteAddr().String())
		fmt.Printf("%s local addr: %s\n", this.name, session.LocalAddr().String())
	} else {
		this.result = err
		fmt.Printf("%s: %s\n", this.name, err.Error())
	} // if

	this.timeout.Done()
}

func (this *testSession) Encode(message any) (packet []byte, err error) {
	if this.encodeFailed {
		return nil, fmt.Errorf("encode failed")
	} // if

	return []byte(message.(string)), nil
}

func (this *testSession) Decode(packet []byte) (message any, err error) {
	if this.decodeFailed {
		return nil, fmt.Errorf("decode failed")
	} // if

	return string(packet), nil
}

func (this *testSession) Start() {
	this.lock.Lock()
	defer this.lock.Unlock()

	fmt.Printf("%s start\n", this.name)
}

func (this *testSession) Finish() {
	this.lock.Lock()
	defer this.lock.Unlock()

	fmt.Printf("%s finish\n", this.name)
}

func (this *testSession) Receive(message any) error {
	if this.name != message {
		return fmt.Errorf("receive failed")
	} // if

	return nil
}

func (this *testSession) Error(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.result = err
}
