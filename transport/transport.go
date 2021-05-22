/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package transport

import (
	"errors"

	"github.com/eric135/YaNFD/ndn"
	"github.com/eric135/YaNFD/ndn/lpv2"
	"github.com/eric135/YaNFD/ndn/tlv"
)

type Transport interface {
	SendAndReceive(interest *ndn.Interest) error
}

type transportBase struct {
}

func (t *transportBase) validateReceivedWire(wire []byte) error {
	block, _, err := tlv.DecodeBlock(wire)
	if err != nil {
		return err
	}

	// Need to decode as lpv2.Packet in case it is in an NDNLPv2 header
	lpPacket, err := lpv2.DecodePacket(block)
	if err != nil {
		return err
	}

	if lpPacket.IsBare() {
		return errors.New("empty response received")
	}

	if len(lpPacket.Fragment()) == 0 || lpPacket.Fragment()[0] != tlv.Data {
		return errors.New("non-data type received")
	}
	return nil
}
