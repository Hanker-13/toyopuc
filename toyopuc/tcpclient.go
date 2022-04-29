package toyopuc

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	// tcpProtocolIdentifier uint16 = 0x0000

	// Modbus Application Protocol
	tcpHeaderSize = 0x04 // FT RC LL LH
	tcpMaxLength  = 0x200
	// Default TCP timeout is not set
	tcpTimeout     = 10 * time.Second
	tcpIdleTimeout = 60 * time.Second
)

// TCPClientHandler implements Packager and Transporter interface.
type TCPClientHandler struct {
	tcpPackager
	tcpTransporter
}

// NewTCPClientHandler allocates a new TCPClientHandler.
func NewTCPClientHandler(address string) *TCPClientHandler {
	h := &TCPClientHandler{}
	h.Address = address
	h.Timeout = tcpTimeout
	h.IdleTimeout = tcpIdleTimeout
	h.RequestFT = RequestFTByte
	h.ResponseFTByte = ResponseFTByte
	return h
}

// TCPClient creates TCP client with default handler and given connect string.
func TCPClient(address string) Client {
	handler := NewTCPClientHandler(address)
	return NewClient(handler)
}

// tcpPackager implements Packager interface.
type tcpPackager struct {
	RequestFT      byte
	ResponseFTByte byte
}

// Encode adds toyopuc application protocol header:
// 编码 组装通信报文
//  FT: 1 bytes 0x00
//  RC: 1 bytes 0x00 // 指令中为占位 统一为0x00
//  LL: 1 bytes // 数据长度 低位
//  LH: 1 bytes // 数据长度 高位
//  CMD: 1 byte // 指令代码
//  Data: n bytes // 数据
func (toyopuc *tcpPackager) Encode(pdu *ProtocolDataUnit) (adu []byte, err error) {
	// 编码指令 CDAB
	adu = make([]byte, tcpHeaderSize+1+len(pdu.Data))
	adu[0] = toyopuc.RequestFT
	adu[1] = 0x00 // 占位 该位表示结果状态
	length := len(pdu.Data) + 1
	binary.LittleEndian.PutUint16(adu[2:], uint16(length))
	adu[tcpHeaderSize] = pdu.FunctionCode
	copy(adu[tcpHeaderSize+1:], pdu.Data)
	return
}

// Verify confirms
// 校验确认
func (toyopuc *tcpPackager) Verify(aduRequest []byte, aduResponse []byte) (err error) {
	// 校验头
	if aduResponse[0] != toyopuc.ResponseFTByte {
		err = fmt.Errorf("toyopuc: FT error")
		return
	}
	// 校验长度
	length := int(binary.LittleEndian.Uint16(aduResponse[2:4]))
	if length != len(aduResponse[tcpHeaderSize:]) {
		err = fmt.Errorf("toyopuc: response data length error")
		return
	}
	return
}

// Decode extracts PDU from TCP frame:
// 解码
//  FT: 1 bytes 0x00
//  RC: 1 bytes // 0x00 正常 ，其他为不正常
//  LL: 1 bytes // 数据长度 低位
//  LH: 1 bytes // 数据长度 高位
//  CMD: 1 byte // 指令代码
//  Data: n bytes // 数据
func (toyopuc *tcpPackager) Decode(adu []byte) (pdu *ProtocolDataUnit, err error) {
	pdu = &ProtocolDataUnit{}
	pdu.FunctionCode = adu[4]
	pdu.Response = adu[1]
	if len(adu) > tcpHeaderSize+1 {
		// 查询结果 CDAB
		pdu.Data = adu[tcpHeaderSize+1:]
	}
	return
}

// tcpTransporter implements Transporter interface.
type tcpTransporter struct {
	// Connect string
	Address string
	// Connect & Read timeout
	Timeout time.Duration
	// Idle timeout to close the connection
	IdleTimeout time.Duration
	// Transmission logger
	Logger *log.Logger

	// TCP connection
	mu           sync.Mutex
	conn         net.Conn
	closeTimer   *time.Timer
	lastActivity time.Time
}

// Send sends data to server and ensures response length is greater than header length.
// 发送
func (toyopuc *tcpTransporter) Send(aduRequest []byte) (aduResponse []byte, err error) {
	toyopuc.mu.Lock()
	defer toyopuc.mu.Unlock()

	// Establish a new connection if not connected
	if err = toyopuc.connect(); err != nil {
		return
	}
	// Set timer to close when idle
	toyopuc.lastActivity = time.Now()
	toyopuc.startCloseTimer()
	// Set write and read timeout
	var timeout time.Time
	if toyopuc.Timeout > 0 {
		timeout = toyopuc.lastActivity.Add(toyopuc.Timeout)
	}
	if err = toyopuc.conn.SetDeadline(timeout); err != nil {
		return
	}
	// Send data
	toyopuc.logf("toyopuc: sending % x", aduRequest)
	if _, err = toyopuc.conn.Write(aduRequest); err != nil {
		return
	}

	// Read header first
	var data [tcpMaxLength]byte
	if _, err = io.ReadFull(toyopuc.conn, data[:tcpHeaderSize]); err != nil {
		return
	}
	// Read length, ignore transaction & protocol id (4 bytes)
	length := int(binary.LittleEndian.Uint16(data[2:]))
	if length <= 0 {
		toyopuc.flush(data[:])
		err = fmt.Errorf("toyopuc: length in response header '%v' must not be zero", length)
		return
	}
	// length := len(data)
	if length > (tcpMaxLength - (tcpHeaderSize - 1)) {
		toyopuc.flush(data[:])
		err = fmt.Errorf("toyopuc: length in response header '%v' must not greater than '%v'", length, tcpMaxLength-tcpHeaderSize+1)
		return
	}
	// Skip unit id
	// 读取不到足够数量 则报timeout
	length += tcpHeaderSize
	if _, err = io.ReadFull(toyopuc.conn, data[tcpHeaderSize:length]); err != nil {
		return
	}
	aduResponse = data[:length]
	toyopuc.logf("toyopuc: received % x\n", aduResponse)
	return
}

// Connect establishes a new connection to the address in Address.
// Connect and Close are exported so that multiple requests can be done with one session
func (toyopuc *tcpTransporter) Connect() error {
	toyopuc.mu.Lock()
	defer toyopuc.mu.Unlock()

	return toyopuc.connect()
}

func (toyopuc *tcpTransporter) connect() error {
	if toyopuc.conn == nil {
		dialer := net.Dialer{Timeout: toyopuc.Timeout}
		conn, err := dialer.Dial("tcp", toyopuc.Address)
		if err != nil {
			return err
		}
		toyopuc.conn = conn
	}
	return nil
}

func (toyopuc *tcpTransporter) startCloseTimer() {
	if toyopuc.IdleTimeout <= 0 {
		return
	}
	if toyopuc.closeTimer == nil {
		toyopuc.closeTimer = time.AfterFunc(toyopuc.IdleTimeout, toyopuc.closeIdle)
	} else {
		toyopuc.closeTimer.Reset(toyopuc.IdleTimeout)
	}
}

// Close closes current connection.
func (toyopuc *tcpTransporter) Close() error {
	toyopuc.mu.Lock()
	defer toyopuc.mu.Unlock()

	return toyopuc.close()
}

// flush flushes pending data in the connection,
// returns io.EOF if connection is closed.
func (toyopuc *tcpTransporter) flush(b []byte) (err error) {
	if err = toyopuc.conn.SetReadDeadline(time.Now()); err != nil {
		return
	}
	// Timeout setting will be reset when reading
	if _, err = toyopuc.conn.Read(b); err != nil {
		// Ignore timeout error
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			err = nil
		}
	}
	return
}

func (toyopuc *tcpTransporter) logf(format string, v ...interface{}) {
	if toyopuc.Logger != nil {
		toyopuc.Logger.Printf(format, v...)
	}
}

// closeLocked closes current connection. Caller must hold the mutex before calling this method.
func (toyopuc *tcpTransporter) close() (err error) {
	if toyopuc.conn != nil {
		err = toyopuc.conn.Close()
		toyopuc.conn = nil
	}
	return
}

// closeIdle closes the connection if last activity is passed behind IdleTimeout.
func (toyopuc *tcpTransporter) closeIdle() {
	toyopuc.mu.Lock()
	defer toyopuc.mu.Unlock()

	if toyopuc.IdleTimeout <= 0 {
		return
	}
	idle := time.Since(toyopuc.lastActivity)
	if idle >= toyopuc.IdleTimeout {
		toyopuc.logf("toyopuc: closing connection due to idle timeout: %v", idle)
		toyopuc.close()
	}
}
