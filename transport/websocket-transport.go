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
	"github.com/eric135/ndn-reachability/uri"
	"github.com/gorilla/websocket"
)

type WebSocketTransport struct {
	conn *websocket.Conn
	transportBase
}

func NewWebSocketTransport(protocol string, wsUri string) (*WebSocketTransport, error) {
	t := new(WebSocketTransport)
	var err error

	var dialer websocket.Dialer
	// Override dialer to ensure we get the specified IP version
	dialer.NetDial = func(network string, addr string) (net.Conn, error) {
		u, err := uri.ParseWebSocket(wsUri)
		if err != nil {
			return nil, err
		}

		if protocol == "wss-ipv4" {
			return net.Dial("tcp4", net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
		} else {
			return net.Dial("tcp6", net.JoinHostPort(u.Host, strconv.Itoa(u.Port)))
		}
	}

	t.conn, _, err = dialer.Dial(wsUri, nil)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *WebSocketTransport) SendAndReceive(interest *ndn.Interest) (time.Duration, error) {
	block, err := interest.Encode()
	if err != nil {
		return 0, errors.New("internal error")
	}
	wire, err := block.Wire()
	if err != nil {
		return 0, errors.New("internal error")
	}
	t.conn.WriteMessage(websocket.BinaryMessage, wire)
	t.conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	startTime := time.Now()

	// Wait for response until timeout
	for time.Now().Before(startTime.Add(time.Second * 5)) {
		messageType, readBuf, err := t.conn.ReadMessage()

		if err != nil {
			// If err not nil, likely timed out or another issue
			return 0, errors.New("timeout")
		}
		if messageType == websocket.CloseMessage {
			// Remote endpoint closed connection
			return 0, errors.New("timeout")
		}

		// Make sure reply is Data packet and not Nack
		err = t.validateReceivedWire(readBuf)
		if err == nil {
			return time.Since(startTime), nil
		}
	}

	// Return last error
	t.conn.Close()
	return 0, err
}
