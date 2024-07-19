package nets

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/yinweli/Mizugo/mizugos/pools"
)

// NewTCPSession 建立TCP會話器
func NewTCPSession(conn net.Conn) *TCPSession {
	return &TCPSession{
		conn:       conn,
		message:    make(chan any, ChannelSize),
		headerSize: HeaderSize,
		packetSize: PacketSize,
	}
}

// TCPSession TCP會話器, 負責傳送/接收訊息等相關的功能
type TCPSession struct {
	conn       net.Conn       // 連接物件
	message    chan any       // 訊息通道
	signal     sync.WaitGroup // 通知信號
	codec      []Codec        // 編碼/解碼
	codecr     []Codec        // 編碼/解碼(反序)
	publish    []Publish      // 發布事件處理
	wrong      []Wrong        // 錯誤處理
	headerSize int            // 標頭長度
	packetSize int            // 封包長度
	owner      any            // 擁有者
}

// Start 啟動會話
func (this *TCPSession) Start(bind Bind, unbind Unbind) {
	pools.DefaultPool.Submit(func() {
		if bind.Do(this) == false {
			unbind.Do(this)
			_ = this.conn.Close() // 綁定失敗就結束了
			return
		} // if

		pools.DefaultPool.Submit(this.recvLoop)
		pools.DefaultPool.Submit(this.sendLoop)
		this.doPublish(EventStart, this)
		this.signal.Add(2)
		this.signal.Wait() // 等待接收循環與傳送循環結束, 如果接收循環與傳送循環結束, 就會繼續進行結束處理
		unbind.Do(this)
		this.doPublish(EventStop, this)
	})
}

// Stop 停止會話, 不會等待會話內部循環結束
func (this *TCPSession) Stop() {
	this.message <- nil // 以空訊息通知會話結束
}

// StopWait 停止會話, 會等待會話內部循環結束
func (this *TCPSession) StopWait() {
	this.message <- nil // 以空訊息通知會話結束
	this.signal.Wait()
}

// SetCodec 設定編碼/解碼
func (this *TCPSession) SetCodec(codec ...Codec) {
	this.codec = codec
	this.codecr = nil

	for i := len(codec) - 1; i >= 0; i-- {
		this.codecr = append(this.codecr, codec[i])
	} // for
}

// SetPublish 設定發布事件處理
func (this *TCPSession) SetPublish(publish ...Publish) {
	this.publish = publish
}

// SetWrong 設定錯誤處理
func (this *TCPSession) SetWrong(wrong ...Wrong) {
	this.wrong = wrong
}

// SetHeaderSize 設定標頭長度
func (this *TCPSession) SetHeaderSize(size int) {
	this.headerSize = size
}

// SetPacketSize 設定封包長度
func (this *TCPSession) SetPacketSize(size int) {
	this.packetSize = size
}

// SetOwner 設定擁有者
func (this *TCPSession) SetOwner(owner any) {
	this.owner = owner
}

// Send 傳送訊息
func (this *TCPSession) Send(message any) {
	if message != nil {
		this.message <- message
	} // if
}

// RemoteAddr 取得遠端位址
func (this *TCPSession) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

// LocalAddr 取得本地位址
func (this *TCPSession) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

// GetOwner 取得擁有者
func (this *TCPSession) GetOwner() any {
	return this.owner
}

// recvLoop 接收循環
func (this *TCPSession) recvLoop() {
	// 由於recvLoop的執行方式, 所以不需要用context方式監控終止方式

	reader := bufio.NewReader(this.conn)

	for {
		packet, err := this.recvPacket(reader)

		if err != nil {
			this.doWrong(fmt.Errorf("tcp session recv loop: %w", err))
			break
		} // if

		message, err := this.doDecode(packet)

		if err != nil {
			this.doWrong(fmt.Errorf("tcp session recv loop: %w", err))
			break
		} // if

		this.doPublish(EventRecv, message)
	} // for

	this.message <- nil // 以空訊息通知會話結束
	this.signal.Done()
}

// recvPacket 接收封包
func (this *TCPSession) recvPacket(reader io.Reader) (packet []byte, err error) {
	header := make([]byte, this.headerSize)

	if _, err = io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("tcp session recv packet: %w", err)
	} // if

	size := binary.LittleEndian.Uint16(header)

	if size <= 0 {
		return []byte{}, nil
	} // if

	packet = make([]byte, size)

	if _, err = io.ReadFull(reader, packet); err != nil {
		return nil, fmt.Errorf("tcp session recv packet: %w", err)
	} // if

	return packet, nil
}

// sendLoop 傳送循環
func (this *TCPSession) sendLoop() {
	// 由於sendLoop的執行方式, 所以不需要用context方式監控終止方式

	for {
		message := <-this.message

		if message == nil { // 空訊息表示會話結束
			break
		} // if

		packet, err := this.doEncode(message)

		if err != nil {
			this.doWrong(fmt.Errorf("tcp session send loop: %w", err))
			break
		} // if

		bytes, ok := packet.([]byte)

		if ok == false {
			this.doWrong(fmt.Errorf("tcp session send loop: encode final output not []byte"))
			break
		} // if

		if err = this.sendPacket(this.conn, bytes); err != nil {
			this.doWrong(fmt.Errorf("tcp session send loop: %w", err))
			break
		} // if

		this.doPublish(EventSend, message)
	} // for

	_ = this.conn.Close()
	this.signal.Done()
}

// sendPacket 傳送封包
func (this *TCPSession) sendPacket(writer io.Writer, packet []byte) (err error) {
	size := len(packet)

	if size <= 0 {
		return nil
	} // if

	if size > this.packetSize {
		return fmt.Errorf("tcp session send packet: packet too large")
	} // if

	header := make([]byte, this.headerSize)
	binary.LittleEndian.PutUint16(header, uint16(size))

	if _, err = writer.Write(header); err != nil {
		return fmt.Errorf("tcp session send packet: %w", err)
	} // if

	if _, err = writer.Write(packet); err != nil {
		return fmt.Errorf("tcp session send packet: %w", err)
	} // if

	return nil
}

// doPublish 執行發布事件處理
func (this *TCPSession) doPublish(name string, param any) {
	for _, itor := range this.publish {
		itor.Do(name, param)
	} // for
}

// doEncode 執行封包編碼處理
func (this *TCPSession) doEncode(input any) (output any, err error) {
	for _, itor := range this.codec {
		if input, err = itor.Encode(input); err != nil {
			return nil, fmt.Errorf("tcp session encode: %w", err)
		} // if
	} // for

	return input, nil
}

// doDecode 執行封包解碼處理
func (this *TCPSession) doDecode(input any) (output any, err error) {
	for _, itor := range this.codecr {
		if input, err = itor.Decode(input); err != nil {
			return nil, fmt.Errorf("tcp session decode: %w", err)
		} // if
	} // for

	return input, nil
}

// doWrong 執行錯誤處理
func (this *TCPSession) doWrong(err error) {
	for _, itor := range this.wrong {
		itor.Do(err)
	} // for
}
