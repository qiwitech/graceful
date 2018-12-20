package graceful

import (
	"context"
	"errors"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

// direct transport

type FakeRoudTripper struct {
	f func(*http.Request) (*http.Response, error)
}

func (f FakeRoudTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.f(req)
}

// listener

type FakeListener struct {
	connections      chan net.Conn
	closed           int32
	DefaultTransport http.RoundTripper
}

func NewFakeListener() (l *FakeListener) {
	l = &FakeListener{
		connections: make(chan net.Conn, 10),
	}
	l.DefaultTransport = l.fakeTransport()
	return l
}

func (l *FakeListener) Addr() net.Addr { return &net.TCPAddr{} }

func (l *FakeListener) Accept() (net.Conn, error) {
	conn, ok := <-l.connections
	if !ok {
		return nil, errors.New("closed fake listener")
	}
	return conn, nil
}

func (l *FakeListener) Close() error {
	if atomic.CompareAndSwapInt32(&l.closed, 0, 1) {
		close(l.connections)
	}
	return nil
}

func (l *FakeListener) fakeTransport() http.RoundTripper {
	return &http.Transport{
		DialContext:           l.dialPipe,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

func (l *FakeListener) dialPipe(ctx context.Context, network, addr string) (net.Conn, error) {
	conn1, conn2 := net.Pipe()
	l.connections <- NewFakeConn(ctx, conn1)
	return NewFakeConn(ctx, conn2), nil
}

type FakeConn struct {
	underConn net.Conn
	ctxErr    atomic.Value
}

func NewFakeConn(ctx context.Context, conn net.Conn) *FakeConn {
	c := &FakeConn{
		underConn: conn,
	}
	go func() {
		<-ctx.Done()
		c.ctxErr.Store(ctx.Err())
	}()
	return c
}

func (c *FakeConn) Read(b []byte) (int, error) {
	err, ok := c.ctxErr.Load().(error)
	if ok && err != nil {
		return 0, err
	}
	return c.underConn.Read(b)
}

func (c *FakeConn) Write(b []byte) (int, error) {
	err, ok := c.ctxErr.Load().(error)
	if ok && err != nil {
		return 0, err
	}
	return c.underConn.Write(b)
}

func (c *FakeConn) Close() error {
	return c.underConn.Close()
}

func (c *FakeConn) LocalAddr() net.Addr {
	return c.underConn.LocalAddr()
}

func (c *FakeConn) RemoteAddr() net.Addr {
	return c.underConn.RemoteAddr()
}

func (c *FakeConn) SetDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *FakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *FakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
