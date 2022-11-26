package nets

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
)

const tcpHeaderSize = 2               // 標頭長度
const tcpPacketSize = int(^uint16(0)) // 封包長度

// NewTCPSession 建立tcp會話器
func NewTCPSession(conn net.Conn) *TCPSession {
	return &TCPSession{
		conn:   conn,
		convey: make(chan []byte, conveySize),
	}
}

// TCPSession tcp會話器
type TCPSession struct {
	conn      net.Conn       // 連接物件
	convey    chan []byte    // 傳送通道
	signal    sync.WaitGroup // 通知信號
	sessionID atomic.Int64   // 會話編號
	receive   Receive        // 接收函式
	inform    Inform         // 通知函式
}

// Start 啟動會話, 當由連接器/監聽器獲得會話器之後, 需要啟動會話才可以傳送或接收封包; 若不是使用多執行緒啟動, 則會被阻塞在這裡直到會話結束
func (this *TCPSession) Start(sessionID SessionID, receive Receive, inform Inform) {
	this.signal.Add(2) // 等待接收循環與傳送循環結束
	this.sessionID.Store(sessionID)
	this.receive = receive
	this.inform = inform

	go this.recvLoop()
	go this.sendLoop()

	this.signal.Wait() // 如果接收循環與傳送循環結束, 就會繼續進行結束處理
	this.inform(nil)
}

// Stop 停止會話, 不會等待會話內部循環結束
func (this *TCPSession) Stop() {
	this.convey <- []byte{} // 以空封包通知會話結束
}

// StopWait 停止會話, 會等待會話內部循環結束
func (this *TCPSession) StopWait() {
	this.convey <- []byte{} // 以空封包通知會話結束
	this.signal.Wait()
}

// Send 傳送封包
func (this *TCPSession) Send(packet []byte) {
	if len(packet) > 0 {
		this.convey <- packet
	} // if
}

// SessionID 取得會話編號
func (this *TCPSession) SessionID() SessionID {
	return this.sessionID.Load()
}

// RemoteAddr 取得遠端位址
func (this *TCPSession) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

// LocalAddr 取得本地位址
func (this *TCPSession) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

// recvLoop 接收循環
func (this *TCPSession) recvLoop() {
	reader := bufio.NewReader(this.conn)

	for {
		packet, err := this.recvPacket(reader)

		if err != nil {
			this.inform(fmt.Errorf("recv loop: %w", err))
			break
		} // if

		this.receive(packet)
	} // for

	this.convey <- nil // 以空封包通知會話結束
	this.signal.Done()
}

// recvPacket 接收封包
func (this *TCPSession) recvPacket(reader io.Reader) (packet []byte, err error) {
	header := make([]byte, tcpHeaderSize)

	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, fmt.Errorf("recv packet: %w", err)
	} // if

	size := binary.LittleEndian.Uint16(header)

	if size <= 0 {
		return []byte{}, nil
	} // if

	packet = make([]byte, size)

	if _, err := io.ReadFull(reader, packet); err != nil {
		return nil, fmt.Errorf("recv packet: %w", err)
	} // if

	return packet, nil
}

// sendLoop 傳送循環
func (this *TCPSession) sendLoop() {
	for {
		packet := <-this.convey

		if len(packet) == 0 { // 空封包表示會話結束
			break
		} // if

		if err := this.sendPacket(this.conn, packet); err != nil {
			this.inform(fmt.Errorf("send loop: %w", err))
			break
		} // if
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
		return fmt.Errorf("send packet: packet too large")
	} // if

	header := make([]byte, tcpHeaderSize)
	binary.LittleEndian.PutUint16(header, uint16(size))

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("send packet: %w", err)
	} // if

	if _, err := writer.Write(packet); err != nil {
		return fmt.Errorf("send packet: %w", err)
	} // if

	return nil
}
