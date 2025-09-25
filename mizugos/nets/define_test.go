package nets

import (
	"fmt"
	"net"
	"sync"

	"github.com/yinweli/Mizugo/v2/mizugos/trials"
)

// newTestNet 建立網路測試器
func newTestNet(bind, encode, decode bool) *testNet {
	trials.WaitTimeout() // 在這邊等待一下, 讓程序有機會完成
	return &testNet{
		bind:   bind,
		encode: encode,
		decode: decode,
	}
}

// testNet 網路測試器
type testNet struct {
	bind        bool
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
	message     any
	session     Sessioner
	lock        sync.RWMutex
}

func (this *testNet) Bind(session Sessioner) bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.bind == false {
		return false
	} // if

	this.bindCount++
	this.session = session
	session.SetCodec(&testCodec{
		encode: func(input any) (output any, err error) {
			this.lock.Lock()
			defer this.lock.Unlock()
			this.encodeCount++

			if this.encode {
				return []byte(input.(string)), nil
			} else {
				return nil, fmt.Errorf("encode failed")
			} // if
		},
		decode: func(input any) (output any, err error) {
			this.lock.Lock()
			defer this.lock.Unlock()
			this.decodeCount++

			if this.decode {
				return string(input.([]byte)), nil
			} else {
				return nil, fmt.Errorf("decode failed")
			} // if
		},
	})
	session.SetPublish(func(name string, param any) {
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
	})
	session.SetWrong(this.Wrong)
	session.SetPacketSize(PacketSize)
	return true
}

func (this *testNet) Unbind(_ Sessioner) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.unbindCount++
	this.session = nil
}

func (this *testNet) Wrong(err error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.err = err
}

func (this *testNet) Valid() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.err == nil
}

func (this *testNet) ValidBind() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.bindCount > 0
}

func (this *testNet) ValidUnbind() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.unbindCount > 0
}

func (this *testNet) ValidEncode() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.encodeCount > 0
}

func (this *testNet) ValidDecode() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.decodeCount > 0
}

func (this *testNet) ValidStart() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.startCount == 1
}

func (this *testNet) ValidStop() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.stopCount == 1
}

func (this *testNet) ValidRecv() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.recvCount > 0
}

func (this *testNet) ValidSend() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.sendCount > 0
}

func (this *testNet) ValidSession() bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.session != nil
}

func (this *testNet) ValidMessage(message any) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.message == message
}

func (this *testNet) Get() Sessioner {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.session
}

// testCodec 測試編碼/解碼
type testCodec struct {
	encode func(input any) (output any, err error)
	decode func(input any) (output any, err error)
}

func (this *testCodec) Encode(input any) (output any, err error) {
	return this.encode(input)
}

func (this *testCodec) Decode(input any) (output any, err error) {
	return this.decode(input)
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

func (this *emptySession) Start(_ Bind, _ Unbind) {
}

func (this *emptySession) Stop() {
}

func (this *emptySession) StopWait() {
}

func (this *emptySession) SetCodec(_ ...Codec) {
}

func (this *emptySession) SetPublish(_ ...Publish) {
}

func (this *emptySession) SetWrong(_ ...Wrong) {
}

func (this *emptySession) SetPacketSize(_ int) {
}

func (this *emptySession) SetOwner(_ any) {
}

func (this *emptySession) GetOwner() any {
	return nil
}

func (this *emptySession) Send(_ any) {
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
