package nets

import (
	"fmt"
	"net"
	"time"

	"github.com/yinweli/Mizugo/mizugo/utils"
)

// 這裡存放用於測試的組件

const testerWaitTime = time.Millisecond * 100

func newSessionTester() *sessionTester {
	time.Sleep(testerWaitTime) // 在這邊等待一下, 讓伺服器有機會完成

	return &sessionTester{
		timeout: utils.NewWaitTimeout(time.Second),
	}
}

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

func newCoderTester(encode, decode bool) *coderTester {
	return &coderTester{
		encode: encode,
		decode: decode,
	}
}

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

func newReactorTester(success bool) *reactorTester {
	return &reactorTester{
		success: success,
	}
}

type reactorTester struct {
	success      bool
	flagActive   bool
	flagInactive bool
	flagError    bool
	flagReceive  bool
	message      any
	err          error
}

func (this *reactorTester) valid(message any) bool {
	return this.message == message && this.err == nil
}

func (this *reactorTester) Active() {
	this.flagActive = true
}

func (this *reactorTester) Inactive() {
	this.flagInactive = true
}

func (this *reactorTester) Error(err error) {
	this.err = err
	this.flagError = true
}

func (this *reactorTester) Receive(message any) error {
	this.message = message
	this.flagReceive = true

	if this.success {
		return nil
	} else {
		return fmt.Errorf("receive failed")
	} // if
}

func newNetmgrTester() *netmgrTester {
	time.Sleep(testerWaitTime) // 在這邊等待一下, 讓伺服器有機會完成

	return &netmgrTester{
		coder:   newCoderTester(true, true),
		reactor: newReactorTester(true),
		timeout: utils.NewWaitTimeout(time.Second),
	}
}

type netmgrTester struct {
	coder   *coderTester
	reactor *reactorTester
	timeout *utils.WaitTimeout
	session Sessioner
}

func (this *netmgrTester) wait() bool {
	return this.timeout.Wait()
}

func (this *netmgrTester) valid() bool {
	return this.session != nil
}

func (this *netmgrTester) get() Sessioner {
	return this.session
}

func (this *netmgrTester) Prepare(session Sessioner) (coder Coder, reactor Reactor) {
	this.session = session
	this.timeout.Done()
	return this.coder, this.reactor
}

func (this *netmgrTester) Error(err error) {
	this.timeout.Done()
}

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
