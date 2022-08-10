package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telClient struct {
	conn    net.Conn
	addr    string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *telClient) Connect() (err error) {
	t.conn, err = net.DialTimeout("tcp", t.addr, t.timeout)
	return err
}

func (t *telClient) Close() error {
	return t.conn.Close()
}

func (t *telClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	return err
}

func (t *telClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telClient{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
