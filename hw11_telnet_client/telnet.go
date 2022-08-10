package main

import (
	"bufio"
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

func (t *telClient) Send() (err error) {
	s := bufio.NewScanner(t.in)
	for s.Scan() {
		_, err = t.conn.Write(append(s.Bytes(), '\n'))
		if err != nil {
			return
		}
	}

	if s.Err() != nil {
		err = s.Err()
	}
	return err
}

func (t *telClient) Receive() (err error) {
	s := bufio.NewScanner(t.conn)
	for s.Scan() {
		_, err = t.out.Write(append(s.Bytes(), '\n'))
		if err != nil {
			return
		}
	}
	if s.Err() != nil {
		err = s.Err()
	}
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
