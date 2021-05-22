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

type TCPTransport struct {
	conn net.Conn
	transportBase
}

func NewTCPTransport(scheme string, host string, port int) *UDPTransport {
	t := new(UDPTransport)
	var err error
	t.conn, err = net.Dial(scheme, host+":"+strconv.Itoa(port))
	if err != nil {
		return nil
	}
	return t
}

func (t *TCPTransport) SendAndReceive(interest *ndn.Interest) error {
	block, err := interest.Encode()
	if err != nil {
		return errors.New("unable to encode Interest")
	}
	wire, err := block.Wire()
	if err != nil {
		return errors.New("unable to encode Interest")
	}
	t.conn.Write(wire)
	t.conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	startTime := time.Now()

	// Wait for response until timeout
	readBuf := make([]byte, 8800)
	for time.Now().Before(startTime.Add(time.Second * 5)) {
		var readSize int
		readSize, err = t.conn.Read(readBuf)

		// If err not nil, likely timed out or another issue
		if err != nil {
			return err
		}

		// Make sure reply is Data packet and not Nack
		err = t.validateReceivedWire(readBuf[:readSize])
		if err == nil {
			return nil
		}
	}

	// Return last error
	t.conn.Close()
	return err
}
