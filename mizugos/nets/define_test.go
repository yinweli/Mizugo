package nets

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/yinweli/Mizugo/testdata"
)

// newInformTester 建立測試器
func newTester(bundle, encode, decode bool) *tester {
	time.Sleep(testdata.Timeout) // 在這邊等待一下, 讓程序有機會完成
	return &tester{
		bundle: bundle,
		encode: encode,
		decode: decode,
	}
}

// tester 測試器
type tester struct {
	bundle      bool
	encode      bool
	decode      bool
	err         error
	bindCount   int
	unbindCount int
	encodeCount int
	decodeCount int
	startCount  int
	stopCount   int
	recvCount   int
	sendCount   int
	session     Sessioner
	message     any
	lock        sync.RWMutex
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

func (this *tester) validStart() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.startCount == 1
}

func (this *tester) validStop() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.stopCount == 1
}

func (this *tester) validRecv() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.recvCount > 0
}

func (this *tester) validSend() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()

	return this.sendCount > 0
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

func (this *tester) bind(session Sessioner) *Bundle {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.bindCount++
	this.session = session

	if this.bundle {
		return &Bundle{
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
			Publish: func(name string, param any) {
				this.lock.Lock()
				defer this.lock.Unlock()

				switch name {
				case EventStart:
					this.startCount++

				case EventStop:
					this.stopCount++

				case EventRecv:
					this.recvCount++
					this.message = param

				case EventSend:
					this.sendCount++
				} // switch
			},
		}
	} else {
		return nil
	} // if
}

func (this *tester) unbind(_ Sessioner) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.unbindCount++
	this.session = nil
}

func (this *tester) wrong(_ bool, err error) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.err = err
}

// emptyConnect 空連接
type emptyConnect struct {
	value int // 為了防止建立空物件時使用到相同的指標, 所以弄個變數來影響他
}

func (this *emptyConnect) Connect(_ Bind, _ Unbind, _ Wrong) {
}

func (this *emptyConnect) Address() string {
	return ""
}

// emptyListen 空接聽
type emptyListen struct {
	value int // 為了防止建立空物件時使用到相同的指標, 所以弄個變數來影響他
}

func (this *emptyListen) Listen(_ Bind, _ Unbind, _ Wrong) {
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

func (this *emptySession) Start(_ Bind, _ Unbind, _ Wrong) {
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

func (this *emptySession) SetOwner(_ any) {
}

func (this *emptySession) GetOwner() any {
	return nil
}

// host 端點資料
type host struct {
	ip   string
	port string
}
