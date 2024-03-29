package sandbox

import (
	"io"
	"net"
	"os"
	"sync/atomic"
	"time"
)

type Socket interface {
	Family() Family

	Type() Socktype

	Fd() uintptr

	File() *os.File

	Close() error

	Bind(addr Sockaddr) error

	Listen(backlog int) error

	Connect(addr Sockaddr) error

	Accept() (Socket, Sockaddr, error)

	Name() (Sockaddr, error)

	Peer() (Sockaddr, error)

	RecvFrom(iovs [][]byte, flags int) (n, rflags int, addr Sockaddr, err error)

	SendTo(iovs [][]byte, addr Sockaddr, flags int) (int, error)

	Shutdown(how int) error

	Error() error

	IsListening() (bool, error)

	IsNonBlock() (bool, error)

	TCPNoDelay() (bool, error)

	RecvBuffer() (int, error)

	SendBuffer() (int, error)

	RecvTimeout() (time.Duration, error)

	SendTimeout() (time.Duration, error)

	SetNonBlock(nonblock bool) error

	SetRecvBuffer(size int) error

	SetSendBuffer(size int) error

	SetRecvTimeout(timeout time.Duration) error

	SetSendTimeout(timeout time.Duration) error

	SetTCPNoDelay(nodelay bool) error

	SetTLSServerName(serverName string) error
}

type Socktype uint8

type Family uint8

func (f Family) String() string {
	switch f {
	case UNIX:
		return "UNIX"
	case INET:
		return "INET"
	case INET6:
		return "INET6"
	default:
		return "UNSPEC"
	}
}

type Protocol uint16

const (
	NOPROTO Protocol = 0
	TCP     Protocol = 6
	UDP     Protocol = 17
)

func (p Protocol) String() string {
	switch p {
	case NOPROTO:
		return "NOPROTO"
	case TCP:
		return "TCP"
	case UDP:
		return "UDP"
	default:
		return "UNKNOWN"
	}
}

type socketConn struct {
	sock  Socket
	laddr net.Addr
	raddr net.Addr
}

func (c *socketConn) Close() error {
	return c.netError("close", c.sock.Close())
}

func (c *socketConn) Read(b []byte) (int, error) {
	n, _, _, err := c.sock.RecvFrom([][]byte{b}, 0)
	if err != nil {
		return 0, c.netError("read", err)
	}
	if n == 0 {
		return 0, io.EOF
	}
	return n, nil
}

func (c *socketConn) ReadFrom(b []byte) (int, net.Addr, error) {
	n, _, sa, err := c.sock.RecvFrom([][]byte{b}, 0)
	if err != nil {
		return -1, nil, c.netError("read", err)
	}
	var addr net.Addr
	switch a := sa.(type) {
	case *SockaddrInet4:
		addr = &net.UDPAddr{IP: a.Addr[:], Port: a.Port}
	case *SockaddrInet6:
		addr = &net.UDPAddr{IP: a.Addr[:], Port: a.Port}
	case *SockaddrUnix:
		addr = &net.UnixAddr{Net: "unixgram", Name: a.Name}
	}
	return n, addr, nil
}

func (c *socketConn) Write(b []byte) (int, error) {
	n, err := c.sock.File().Write(b)
	return n, c.netError("write", err)
}

func (c *socketConn) WriteTo(b []byte, addr net.Addr) (int, error) {
	var sa Sockaddr
	switch a := addr.(type) {
	case *net.UDPAddr:
		if ipv4 := a.IP.To4(); ipv4 != nil {
			sa = &SockaddrInet4{Addr: ([4]byte)(ipv4), Port: a.Port}
		} else {
			sa = &SockaddrInet6{Addr: ([16]byte)(a.IP), Port: a.Port}
		}
	case *net.UnixAddr:
		sa = &SockaddrUnix{Name: a.Name}
	}
	n, err := c.sock.SendTo([][]byte{b}, sa, 0)
	return n, c.netError("write", err)
}

func (c *socketConn) LocalAddr() net.Addr {
	return c.laddr
}

func (c *socketConn) RemoteAddr() net.Addr {
	return c.raddr
}

func (c *socketConn) SetDeadline(t time.Time) error {
	return c.netError("set", c.sock.File().SetDeadline(t))
}

func (c *socketConn) SetReadDeadline(t time.Time) error {
	return c.netError("set", c.sock.File().SetReadDeadline(t))
}

func (c *socketConn) SetWriteDeadline(t time.Time) error {
	return c.netError("set", c.sock.File().SetWriteDeadline(t))
}

func (c *socketConn) netError(op string, err error) error {
	return netError(op, c.LocalAddr(), c.RemoteAddr(), err)
}

type socketListener struct {
	sock Socket
	addr net.Addr
}

func (l *socketListener) Close() error {
	return l.netError("close", l.sock.Close())
}

func (l *socketListener) Addr() net.Addr {
	return l.addr
}

func (l *socketListener) Accept() (net.Conn, error) {
	sock, addr, err := l.sock.Accept()
	if err != nil {
		return nil, l.netError("accept", err)
	}
	conn := &socketConn{
		sock:  sock,
		laddr: l.addr,
	}
	switch a := addr.(type) {
	case *SockaddrInet4:
		conn.raddr = &net.TCPAddr{IP: a.Addr[:], Port: a.Port}
	case *SockaddrInet6:
		conn.raddr = &net.TCPAddr{IP: a.Addr[:], Port: a.Port}
	}
	return conn, nil
}

func (l *socketListener) netError(op string, err error) error {
	return netError(op, l.Addr(), nil, err)
}

func netError(op string, laddr, raddr net.Addr, err error) error {
	if err == nil {
		return nil
	}
	if err == io.EOF {
		return err
	}
	return &net.OpError{
		Op:     op,
		Net:    laddr.Network(),
		Source: laddr,
		Addr:   raddr,
		Err:    err,
	}
}

func setFileDeadline(f *os.File, rtimeout, wtimeout time.Duration) error {
	var now time.Time
	if rtimeout > 0 || wtimeout > 0 {
		now = time.Now()
	}
	if rtimeout > 0 {
		if err := f.SetReadDeadline(now.Add(rtimeout)); err != nil {
			return err
		}
	}
	if wtimeout > 0 {
		if err := f.SetWriteDeadline(now.Add(wtimeout)); err != nil {
			return err
		}
	}
	return nil
}

func handleSocketIOError(err error) error {
	if err != nil {
		if err == os.ErrDeadlineExceeded {
			err = EAGAIN
		}
	}
	return err
}

// socketFD is used to manage the lifecycle of socket file descriptors;
// it allows multiple goroutines to share ownership of the socket while
// coordinating to close the file descriptor via an atomic reference count.
//
// Goroutines must call acquire to access the file descriptor; if they get a
// negative number, it indicates that the socket was already closed and the
// method should usually return EBADF.
//
// After acquiring a valid file descriptor, the goroutine is responsible for
// calling release with the same fd number that was returned by acquire. The
// release may cause the file descriptor to be closed if the close method was
// called in between and releasing the fd causes the reference count to reach
// zero.
//
// The close method detaches the file descriptor from the socketFD, but it only
// closes it if the reference count is zero (no other goroutines was sharing
// ownership). After closing the socketFD, all future calls to acquire return a
// negative number, preventing other goroutines from acquiring ownership of the
// file descriptor and guaranteeing that it will eventually be closed.
type socketFD struct {
	state atomic.Uint64 // upper 32 bits: refCount, lower 32 bits: fd
}

func (s *socketFD) init(fd int) {
	s.state.Store(uint64(fd & 0xFFFFFFFF))
}

func (s *socketFD) load() int {
	return int(int32(s.state.Load()))
}

func (s *socketFD) refCount() int {
	return int(s.state.Load() >> 32)
}

func (s *socketFD) acquire() int {
	for {
		oldState := s.state.Load()
		refCount := (oldState >> 32) + 1
		newState := (refCount << 32) | (oldState & 0xFFFFFFFF)

		fd := int32(oldState)
		if fd < 0 {
			return -1
		}
		if s.state.CompareAndSwap(oldState, newState) {
			return int(fd)
		}
	}
}

func (s *socketFD) releaseFunc(fd int, closeFD func(int)) {
	for {
		oldState := s.state.Load()
		refCount := (oldState >> 32) - 1
		newState := (refCount << 32) | (oldState & 0xFFFFFFFF)

		if s.state.CompareAndSwap(oldState, newState) {
			if int32(oldState) < 0 && refCount == 0 {
				closeFD(fd)
			}
			break
		}
	}
}

func (s *socketFD) closeFunc(closeFD func(int)) {
	for {
		oldState := s.state.Load()
		refCount := oldState >> 32
		newState := oldState | 0xFFFFFFFF

		fd := int32(oldState)
		if fd < 0 {
			break
		}
		if s.state.CompareAndSwap(oldState, newState) {
			if refCount == 0 {
				closeFD(int(fd))
			}
			break
		}
	}
}
