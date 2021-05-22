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
	transportBase
}

func NewUDPTransport(scheme string, host string, port int) (*UDPTransport, error) {
	t := new(UDPTransport)
	var err error
	t.conn, err = net.Dial(scheme, host+":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *UDPTransport) SendAndReceive(interest *ndn.Interest) (time.Duration, error) {
	block, err := interest.Encode()
	if err != nil {
		return 0, errors.New("internal error")
	}
	wire, err := block.Wire()
	if err != nil {
		return 0, errors.New("internal error")
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
			return 0, errors.New("timeout")
		}

		// Make sure reply is Data packet and not Nack
		err = t.validateReceivedWire(readBuf[:readSize])
		if err == nil {
			return time.Since(startTime), nil
		}
	}

	// Return last error
	t.conn.Close()
	return 0, err
}
