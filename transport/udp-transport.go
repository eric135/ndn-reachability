/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package transport

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/eric135/YaNFD/ndn"
)

type UDPTransport struct {
	conn net.Conn
}

func NewUDPTransport(scheme string, host string, port int) *UDPTransport {
	t := new(UDPTransport)
	var err error
	t.conn, err = net.Dial(scheme, host+":"+strconv.Itoa(port))
	if err != nil {
		return nil
	}
	return t
}

func (t *UDPTransport) SendAndReceive(interest *ndn.Interest) error {
	block, err := interest.Encode()
	if err != nil {
		return errors.New("unable to encode Interest")
	}
	wire, err := block.Wire()
	if err != nil {
		return errors.New("unable to encode Interest")
	}
	t.conn.Write(wire)

	// Wait for response until timeout
	t.conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	readBuf := make([]byte, 8800)
	_, err = t.conn.Read(readBuf)
	t.conn.Close()

	// If err not nil, likely timed out or another issue
	return err
}
