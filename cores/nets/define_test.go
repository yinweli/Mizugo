package nets

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/testdata"
)

// newDoneTester 建立完成會話測試器
func newDoneTester() *doneTester {
	time.Sleep(testdata.Timeout) // 在這邊等待一下, 讓程序有機會完成
	return &doneTester{}
}

// doneTester 完成會話測試器
type doneTester struct {
	lock    sync.Mutex
	session Sessioner
	err     error
}

func (this *doneTester) valid() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session != nil && this.err == nil
}

func (this *doneTester) get() Sessioner {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session
}

func (this *doneTester) done(session Sessioner, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = session
	this.err = err
}

// newBindTester 建立綁定處理測試器
func newBindTester(encode, decode, receive bool) *bindTester {
	time.Sleep(testdata.Timeout) // 在這邊等待一下, 讓程序有機會完成
	return &bindTester{
		encode:  encode,
		decode:  decode,
		receive: receive,
	}
}

// bindTester 綁定處理測試器
type bindTester struct {
	encode  bool
	decode  bool
	receive bool
	lock    sync.Mutex
	session Sessioner
	message any
	err     error
}

func (this *bindTester) validSession() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session != nil
}

func (this *bindTester) validMessage(message any) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.message == message
}

func (this *bindTester) validError() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.err == nil
}

func (this *bindTester) get() Sessioner {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session
}

func (this *bindTester) Bind(session Sessioner) (content Content, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.session = session

	return Content{
		Unbind: func() {
			this.lock.Lock()
			defer this.lock.Unlock()

			this.session = nil
			this.message = nil
		},
		Encode: func(message any) (packet []byte, err error) {
			if this.encode == false {
				return nil, fmt.Errorf("encode failed")
			} // if

			return []byte(message.(string)), nil
		},
		Decode: func(packet []byte) (message any, err error) {
			if this.decode == false {
				return nil, fmt.Errorf("decode failed")
			} // if

			return string(packet), nil
		},
		Receive: func(message any) error {
			this.lock.Lock()
			defer this.lock.Unlock()

			if this.receive {
				this.message = message
				return nil
			} else {
				this.message = nil
				return fmt.Errorf("failed")
			} // if
		},
	}, nil
}

func (this *bindTester) Error(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.err = err
}

// emptySession 空會話
type emptySession struct {
}

func (this *emptySession) Start(_ SessionID, _ Binder) {
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

// host 端點資料
type host struct {
	ip   string
	port string
}
