package nets

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/testdata"
)

// newInformTester 建立測試器
func newTester(encode, decode, receive bool) *tester {
	time.Sleep(testdata.Timeout) // 在這邊等待一下, 讓程序有機會完成
	return &tester{
		encode:  encode,
		decode:  decode,
		receive: receive,
	}
}

// tester 測試器
type tester struct {
	encode         bool
	decode         bool
	receive        bool
	err            error
	bindCount      int
	unbindCount    int
	encodeCount    int
	decodeCount    int
	receiveCount   int
	afterSendCount int
	afterRecvCount int
	session        Sessioner
	message        any
	lock           sync.RWMutex
}

func (this *tester) valid() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.err == nil
}

func (this *tester) validBind() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.bindCount > 0
}

func (this *tester) validUnbind() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.unbindCount > 0
}

func (this *tester) validEncode() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.encodeCount > 0
}

func (this *tester) validDecode() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.decodeCount > 0
}

func (this *tester) validReceive() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.receiveCount > 0
}

func (this *tester) validAfterSend() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.afterSendCount > 0
}

func (this *tester) validAfterRecv() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.afterRecvCount > 0
}

func (this *tester) validSession() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.session != nil
}

func (this *tester) validMessage(message any) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.message == message
}

func (this *tester) get() Sessioner {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.session
}

func (this *tester) inform() Inform {
	return Inform{
		Error: func(err error) {
			this.lock.Lock()
			defer this.lock.Unlock()

			this.err = err
		},
		Bind: func(session Sessioner) Bundle {
			this.lock.Lock()
			defer this.lock.Unlock()

			this.bindCount++
			this.session = session
			return Bundle{
				Encode: func(message any) (packet []byte, err error) {
					this.lock.Lock()
					defer this.lock.Unlock()

					this.encodeCount++

					if this.encode {
						return []byte(message.(string)), nil
					} else {
						return nil, fmt.Errorf("encode failed")
					} // if
				},
				Decode: func(packet []byte) (message any, err error) {
					this.lock.Lock()
					defer this.lock.Unlock()

					this.decodeCount++

					if this.decode {
						return string(packet), nil
					} else {
						return nil, fmt.Errorf("decode failed")
					} // if
				},
				Receive: func(message any) error {
					this.lock.Lock()
					defer this.lock.Unlock()

					this.receiveCount++

					if this.receive {
						this.message = message
						return nil
					} else {
						this.message = nil
						return fmt.Errorf("failed")
					} // if
				},
				AfterSend: func() {
					this.lock.Lock()
					defer this.lock.Unlock()

					this.afterSendCount++
				},
				AfterRecv: func() {
					this.lock.Lock()
					defer this.lock.Unlock()

					this.afterRecvCount++
				},
			}
		},
		Unbind: func(_ Sessioner) {
			this.lock.Lock()
			defer this.lock.Unlock()

			this.unbindCount++
			this.session = nil
		},
	}
}

// emptyConnect 空連接
type emptyConnect struct {
	value int // 為了防止建立空物件時使用到相同的指標, 所以弄個變數來影響他
}

func (this *emptyConnect) Connect(_ Inform) {
}

func (this *emptyConnect) Address() string {
	return ""
}

// emptyListen 空接聽
type emptyListen struct {
	value int // 為了防止建立空物件時使用到相同的指標, 所以弄個變數來影響他
}

func (this *emptyListen) Listen(_ Inform) {
}

func (this *emptyListen) Stop() error {
	return nil
}

func (this *emptyListen) Address() string {
	return ""
}

// emptySession 空會話
type emptySession struct {
	value int // 為了防止建立空物件時使用到相同的指標, 所以弄個變數來影響他
}

func (this *emptySession) Start(_ Inform) {
}

func (this *emptySession) Stop() {
}

func (this *emptySession) StopWait() {
}

func (this *emptySession) Send(_ any) {
}

func (this *emptySession) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (this *emptySession) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (this *emptySession) SetOwner(owner any) {
}

func (this *emptySession) GetOwner() any {
	return nil
}

// host 端點資料
type host struct {
	ip   string
	port string
}
