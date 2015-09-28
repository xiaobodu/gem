package game

import (
	"net"
	"time"

	"github.com/gtank/isaac"

	"gem/encoding"
	"gem/log"
	"gem/protocol"
)

const (
	GameService   encoding.Int8 = 14
	UpdateService encoding.Int8 = 15
)

// decodeFunc is the function currently used for parsing the read stream and
// dealing with the incoming data.
// If an error is returned, it is assumed that we didn't have enough data, and
// the underlying buffer's read pointer is not altered.
type decodeFunc func(*context, *encoding.Buffer) error

// GameConnection is a network-level representation of the connection.
// It handles read/write buffering, and decodes data into game packets or update requests for processing
type GameConnection struct {
	Index Index
	Log   *log.Module

	randIn  isaac.ISAAC
	randOut isaac.ISAAC
	randKey [4]int32

	conn        net.Conn
	readBuffer  *encoding.Buffer
	writeBuffer *encoding.Buffer
	decode      decodeFunc
	disconnect  chan int
	canRead     chan int
	canWrite    chan int
	active      bool
}

func newConnection(index Index, conn net.Conn, parentLogger *log.Module) *GameConnection {
	return &GameConnection{
		Log:   parentLogger.SubModule(conn.RemoteAddr().String()),
		Index: index,

		conn:        conn,
		readBuffer:  encoding.NewBuffer(),
		writeBuffer: encoding.NewBuffer(),
		disconnect:  make(chan int),
		canRead:     make(chan int, 2),
		canWrite:    make(chan int, 2),
		active:      true,
	}
}

// disconnect signals to the connection loop that this connection should be, or has been closed
func (conn *GameConnection) Disconnect() {
	conn.active = false
	conn.disconnect <- 1
}

// handshake reads the service selection byte and points the connection's decode func
// towards the decode func for the selected service
func (conn *GameConnection) handshake(ctx *context, b *encoding.Buffer) error {
	var svc protocol.ServiceSelect
	if err := svc.Decode(b, nil); err != nil {
		return err
	}

	switch svc.Service {
	case UpdateService:
		if err := new(protocol.UpdateHandshakeResponse).Encode(conn, nil); err != nil {
			return err
		}
		conn.canWrite <- 1

		conn.Log.Infof("new update client")
		conn.decode = ctx.update.handleUpdateRequest
		return nil
	case GameService:
		conn.Log.Infof("new game client")
		conn.decode = ctx.game.handshake
		return nil
	default:
		conn.Log.Errorf("invalid service requested: %v", svc)
		conn.Disconnect()
	}

	return nil
}

// Write is a convenience wrapper around writeBuffer.Write(p)
func (conn *GameConnection) Write(p []byte) (n int, err error) {
	return conn.writeBuffer.Write(p)
}

// flushWriteBuffer drains the write buffer and ensures that all data is written to
// the connection. If conn.Write returns an error (timeout), the client is disconnected.
func (conn *GameConnection) flushWriteBuffer() {
	for conn.writeBuffer.Len() > 0 {
		_, err := conn.writeBuffer.WriteTo(conn.conn)
		if err != nil {
			conn.Log.Debug("write error")
			conn.Disconnect()
			break
		}
	}
	conn.writeBuffer.Trim()
}

// fillReadBuffer pulls data from the connection and buffers it for decoding into the readBuffer
// launched in a goroutine by Server.handle
func (conn *GameConnection) fillReadBuffer() {
	for {
		conn.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err := conn.readBuffer.ReadFrom(conn.conn)
		if err != nil {
			conn.Log.Debugf("read error: %v", err)
			conn.Disconnect()
			break
		}

		conn.canRead <- 1
	}
}