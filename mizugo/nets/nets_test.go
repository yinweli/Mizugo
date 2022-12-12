package nets

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// newCompleteTester 建立完成會話測試器
func newCompleteTester() *completeTester {
	time.Sleep(time.Second) // 在這邊等待一下, 讓程序有機會完成
	return &completeTester{}
}

// completeTester 完成會話測試器
type completeTester struct {
	lock    sync.Mutex
	session Sessioner
	err     error
}

func (this *completeTester) valid() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session != nil && this.err == nil
}

func (this *completeTester) get() Sessioner {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session
}

func (this *completeTester) Complete(session Sessioner, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = session
	this.err = err
}

// newSessionTester 建立會話測試器
func newSessionTester(encode, decode, receive bool) *sessionTester {
	time.Sleep(time.Second) // 在這邊等待一下, 讓程序有機會完成
	return &sessionTester{
		encode:  encode,
		decode:  decode,
		receive: receive,
	}
}

// sessionTester 會話測試器
type sessionTester struct {
	encode  bool
	decode  bool
	receive bool
	lock    sync.Mutex
	session Sessioner
	message any
	err     error
}

func (this *sessionTester) validMessage(message any) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.message == message
}

func (this *sessionTester) validError() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.err == nil
}

func (this *sessionTester) Bind(session Sessioner) (releaser Releaser, reactor Reactor) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = session
	return this, this
}

func (this *sessionTester) Release() {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = nil
	this.message = nil
}

func (this *sessionTester) Encode(message any) (packet []byte, err error) {
	if this.encode == false {
		return nil, fmt.Errorf("encode failed")
	} // if

	return []byte(message.(string)), nil
}

func (this *sessionTester) Decode(packet []byte) (message any, err error) {
	if this.decode == false {
		return nil, fmt.Errorf("decode failed")
	} // if

	return string(packet), nil
}

func (this *sessionTester) Receive(message any) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.receive {
		this.message = message
		return nil
	} else {
		this.message = nil
		return fmt.Errorf("receive failed")
	} // if
}

func (this *sessionTester) Error(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.err = err
}

// emptySession 空會話
type emptySession struct {
}

func (this *emptySession) StartStart(_ SessionID, _ Bundler) {
	// do nothing...
}

func (this *emptySession) Stop() {
	// do nothing...
}

func (this *emptySession) StopWait() {
	// do nothing...
}

func (this *emptySession) Send(_ any) {
	// do nothing...
}

func (this *emptySession) SessionID() SessionID {
	return 0
}

func (this *emptySession) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (this *emptySession) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}
