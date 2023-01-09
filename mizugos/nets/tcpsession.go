package nets

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

const tcpHeaderSize = 2               // 標頭長度
const tcpPacketSize = int(^uint16(0)) // 封包長度
const tcpMessageSize = 1000           // 訊息通道大小設為1000, 避免因為爆滿而卡住

// NewTCPSession 建立tcp會話器
func NewTCPSession(conn net.Conn) *TCPSession {
	return &TCPSession{
		conn:    conn,
		message: make(chan any, tcpMessageSize),
	}
}

// TCPSession tcp會話器
type TCPSession struct {
	conn    net.Conn       // 連接物件
	message chan any       // 訊息通道
	signal  sync.WaitGroup // 通知信號
	bundle  Bundle         // 綁定資料
	wrong   Wrong          // 錯誤處理
	owner   any            // 擁有者
}

// Start 啟動會話
func (this *TCPSession) Start(bind Bind, unbind Unbind, wrong Wrong) {
	go func() {
		this.bundle = bind.Do(this)
		this.wrong = wrong
		this.signal.Add(2) // 等待接收循環與傳送循環結束

		go this.recvLoop()
		go this.sendLoop()

		this.signal.Wait() // 如果接收循環與傳送循環結束, 就會繼續進行結束處理
		unbind.Do(this)
	}()
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

// SetOwner 設定擁有者
func (this *TCPSession) SetOwner(owner any) {
	this.owner = owner
}

// GetOwner 取得擁有者
func (this *TCPSession) GetOwner() any {
	return this.owner
}

// recvLoop 接收循環
func (this *TCPSession) recvLoop() {
	reader := bufio.NewReader(this.conn)

	for {
		packet, err := this.recvPacket(reader)

		if err != nil {
			this.wrong.Do(fmt.Errorf("tcp session recv loop: %w", err))
			break
		} // if

		message, err := this.bundle.Decode(packet)

		if err != nil {
			this.wrong.Do(fmt.Errorf("tcp session recv loop: %w", err))
			break
		} // if

		if err := this.bundle.Receive(message); err != nil {
			this.wrong.Do(fmt.Errorf("tcp session recv loop: %w", err))
			break
		} // if

		this.bundle.AfterRecv.Do()
	} // for

	this.message <- nil // 以空訊息通知會話結束
	this.signal.Done()
}

// recvPacket 接收封包
func (this *TCPSession) recvPacket(reader io.Reader) (packet []byte, err error) {
	header := make([]byte, tcpHeaderSize)

	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("tcp session recv packet: %w", err)
	} // if

	size := binary.LittleEndian.Uint16(header)

	if size <= 0 {
		return []byte{}, nil
	} // if

	packet = make([]byte, size)

	if _, err := io.ReadFull(reader, packet); err != nil {
		return nil, fmt.Errorf("tcp session recv packet: %w", err)
	} // if

	return packet, nil
}

// sendLoop 傳送循環
func (this *TCPSession) sendLoop() {
	for {
		message := <-this.message

		if message == nil { // 空訊息表示會話結束
			break
		} // if

		packet, err := this.bundle.Encode(message)

		if err != nil {
			this.wrong.Do(fmt.Errorf("tcp session send loop: %w", err))
			break
		} // if

		if err := this.sendPacket(this.conn, packet); err != nil {
			this.wrong.Do(fmt.Errorf("tcp session send loop: %w", err))
			break
		} // if

		this.bundle.AfterSend.Do()
	} // for

	_ = this.conn.Close()
	this.signal.Done()
}

// sendPacket 傳送封包
func (this *TCPSession) sendPacket(writer io.Writer, packet []byte) error {
	size := len(packet)

	if size <= 0 {
		return nil
	} // if

	if size > tcpPacketSize {
		return fmt.Errorf("tcp session send packet: packet too large")
	} // if

	header := make([]byte, tcpHeaderSize)
	binary.LittleEndian.PutUint16(header, uint16(size))

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("tcp session send packet: %w", err)
	} // if

	if _, err := writer.Write(packet); err != nil {
		return fmt.Errorf("tcp session send packet: %w", err)
	} // if

	return nil
}
