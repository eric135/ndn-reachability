/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package main

import (
	"fmt"
	"os"

	"github.com/eric135/YaNFD/ndn"
	"github.com/eric135/ndn-reachability/transport"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:")
		fmt.Println(os.Args[0] + " [face-uri] [prefix]")
		os.Exit(-1)
	}

	uri, err := ParseURI(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: Unable to decode Face URI:", err)
	}

	prefix, err := ndn.NameFromString(os.Args[2])
	if err != nil {
		fmt.Println("ERROR: Incorrect name prefix")
		os.Exit(-2)
	}

	fmt.Println("Sending Interest for", prefix.String(), "to"+os.Args[1])

	var t transport.Transport
	if uri.scheme == "udp" || uri.scheme == "udp4" || uri.scheme == "udp6" {
		t = transport.NewUDPTransport(uri.scheme, uri.host, uri.port)
	} else if uri.scheme == "tcp" || uri.scheme == "tcp4" || uri.scheme == "tcp6" {
		// TODO
	} else if uri.scheme == "ws" {
		// TODO
	}

	// Make interest
	interest := ndn.NewInterest(prefix)

	// Send interest and get whether response received
	if t.SendAndReceive(interest) != nil {
		fmt.Println("FAILED: Did not receive Data reply")
		os.Exit(1)
	}

	fmt.Println("SUCCESS: Received Data reply")
}
