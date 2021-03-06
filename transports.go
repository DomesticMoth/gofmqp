/*
    This file is part of gofmqp.

    gofmqp is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    gofmqp is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with gofmqp.  If not, see <https://www.gnu.org/licenses/>.
*/
package gofmqp

import (
	"io"
	"net"
	"strings"
	"errors"
	"crypto/tls"
)


type Connection struct{
	reader MsgReader
	writer MsgWriter
	conn io.Closer
}

func (c *Connection) Close() error {
	return c.Close()
}

func (c *Connection) NextRaw() (msg RawMessage, err error) {
	return c.reader.NextRaw()
}

func (c *Connection) NextUnchecked() (msg Message, err error) {
	return c.reader.NextUnchecked()
}

func (c *Connection) Next() (msg Message, err error) {
	return c.reader.Next()
}

func (c *Connection) SendRaw(msg *RawMessage) (err error) {
	return c.writer.SendRaw(msg)
}

func (c *Connection) Send(msg *Message) (err error) {
	return c.writer.Send(msg)
}

type Listener struct {
	l net.Listener
}

func (l *Listener) Accept() (conn Connection, err error) {
	c, err :=  l.l.Accept()
	if err != nil { return }
	return Connection{NewMsgReader(c), NewMsgWriter(c), c}, err
}

func (l *Listener) Close() error {
	return l.l.Close()
}

func Dial(network, address string, config *tls.Config) (conn Connection, err error) {
	if network == "tcp" || network == "unix" {
		c, err := net.Dial(network, address)
		if err != nil { return conn, err }
		return Connection{NewMsgReader(c), NewMsgWriter(c), c}, err
	}
	if network == "tls" {
		c, err := tls.Dial(network, address, config)
		if err != nil { return conn, err  }
		return Connection{NewMsgReader(c), NewMsgWriter(c), c}, err
	}
	if strings.HasPrefix(network, "tls-") {
		network = strings.TrimPrefix(network, "tls-")
		c, err := tls.Dial(network, address, config)
		if err != nil { return conn, err  }
		return Connection{NewMsgReader(c), NewMsgWriter(c), c}, err
	}
	if strings.HasPrefix(network, "tls+") {
		network = strings.TrimPrefix(network, "tls+")
		c, err := tls.Dial(network, address, config)
		if err != nil { return conn, err  }
		return Connection{NewMsgReader(c), NewMsgWriter(c), c}, err
	}
	return conn, errors.New("Unknown transport type: "+network)
}

func Listen(network, laddr string, config *tls.Config) (listener Listener, err error) {
	if network == "tcp" || network == "unix" {
		l, err := net.Listen(network, laddr)
		if err != nil { return listener, err }
		return Listener{l}, err
	}
	if network == "tls" {
		l, err := tls.Listen(network, laddr, config)
		if err != nil { return listener, err }
		return Listener{l}, err
	}
	if strings.HasPrefix(network, "tls-") {
		network = strings.TrimPrefix(network, "tls-")
		l, err := tls.Listen(network, laddr, config)
		if err != nil { return listener, err }
		return Listener{l}, err
	}
	if strings.HasPrefix(network, "tls+") {
		network = strings.TrimPrefix(network, "tls+")
		l, err := tls.Listen(network, laddr, config)
		if err != nil { return listener, err }
		return Listener{l}, err
	}
	return listener, errors.New("Unknown transport type: "+network)
}
