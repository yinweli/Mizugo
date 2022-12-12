package nets

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/mizugo/utils"
)

// 這裡存放用於單元測試的組件

const testerWaitTime = time.Millisecond * 100

// newSessionTester 建立會話測試器
func newSessionTester() *sessionTester {
	time.Sleep(testerWaitTime) // 在這邊等待一下, 讓伺服器有機會完成

	return &sessionTester{
		timeout: utils.NewWaitTimeout(time.Second),
	}
}

// sessionTester 會話測試器
type sessionTester struct {
	timeout *utils.WaitTimeout
	session Sessioner
	err     error
}

func (this *sessionTester) wait() bool {
	return this.timeout.Wait()
}

func (this *sessionTester) valid() bool {
	return this.session != nil && this.err == nil
}

func (this *sessionTester) get() Sessioner {
	return this.session
}

func (this *sessionTester) complete(session Sessioner, err error) {
	this.session = session
	this.err = err
	this.timeout.Done()
}

// newCoderTester 建立編碼測試器
func newCoderTester(encode, decode bool) *coderTester {
	return &coderTester{
		encode: encode,
		decode: decode,
	}
}

// coderTester 編碼測試器
type coderTester struct {
	encode bool
	decode bool
}

func (this *coderTester) Encode(message any) (packet []byte, err error) {
	if this.encode == false {
		return nil, fmt.Errorf("encode failed")
	} // if

	return []byte(message.(string)), nil
}

func (this *coderTester) Decode(packet []byte) (message any, err error) {
	if this.decode == false {
		return nil, fmt.Errorf("decode failed")
	} // if

	return string(packet), nil
}

// newReactorTester 建立反應測試器
func newReactorTester(receive bool) *reactorTester {
	return &reactorTester{
		receive: receive,
	}
}

// reactorTester 反應測試器
type reactorTester struct {
	receive      bool
	flagActive   bool
	flagInactive bool
	flagError    bool
	flagReceive  bool
	session      Sessioner
	message      any
}

func (this *reactorTester) validSession() bool {
	return this.session != nil
}

func (this *reactorTester) validMessage(message any) bool {
	return this.message == message
}

func (this *reactorTester) get() Sessioner {
	return this.session
}

func (this *reactorTester) Active(session Sessioner) {
	this.flagActive = true
	this.session = session
}

func (this *reactorTester) Inactive() {
	this.flagInactive = true
}

func (this *reactorTester) Error(err error) {
	this.flagError = true
}

func (this *reactorTester) Receive(message any) error {
	this.flagReceive = true

	if this.receive {
		this.message = message
		return nil
	} else {
		this.message = nil
		return fmt.Errorf("receive failed")
	} // if
}

// newNetmgrTester 建立網路管理測試器
func newNetmgrTester() *netmgrTester {
	time.Sleep(testerWaitTime) // 在這邊等待一下, 讓伺服器有機會完成

	return &netmgrTester{
		coder:   newCoderTester(true, true),
		reactor: newReactorTester(true),
	}
}

// netmgrTester 網路管理測試器
type netmgrTester struct {
	coder   *coderTester
	reactor *reactorTester
	success bool
}

func (this *netmgrTester) valid() bool {
	return this.reactor.validSession() && this.success
}

func (this *netmgrTester) sessionID() SessionID {
	return this.reactor.get().SessionID()
}

func (this *netmgrTester) Create() (coder Coder, reactor Reactor) {
	this.success = true
	return this.coder, this.reactor
}

func (this *netmgrTester) Failed(_ error) {
	this.success = false
}

// emptySession 空會話
type emptySession struct {
}

func (this *emptySession) Start(_ SessionID, _ Coder, _ Reactor) {
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
