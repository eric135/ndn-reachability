/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package main

import (
	"errors"
	"strconv"
	"strings"
)

type URI struct {
	scheme string
	host   string
	port   int
	path   string
}

func parseUDP(str string) (*URI, error) {
	uri := new(URI)
	uri.scheme = "udp"
	split := strings.SplitN(str, ":", 2)
	if len(split) < 2 {
		return nil, errors.New("incorrect URI format")
	}
	uri.host = split[0]
	var err error
	uri.port, err = strconv.Atoi(split[1])
	if err != nil {
		return nil, errors.New("could not convert port to integer")
	}
	if uri.port <= 0 || uri.port > 65535 {
		return nil, errors.New("port out of range")
	}
	return uri, nil
}

func parseWebSocket(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 4)
	if len(split1) < 4 {
		return nil, errors.New("incorrect URI format")
	}
	uri.scheme = split1[0][:len(split1[0])-1] // Remove trailing :
	if uri.scheme != "ws" && uri.scheme != "wss" {
		return nil, errors.New("unknown scheme")
	}
	uri.path = split1[3]

	// split1[1] is part between // after scheme
	split2 := strings.SplitN(split1[2], ":", 2)
	var err error
	if len(split2) == 2 {
		// Has port specified
		uri.host = split2[0]
		uri.port, err = strconv.Atoi(split2[1])
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
		if uri.port <= 0 || uri.port > 65535 {
			return nil, errors.New("port out of range")
		}
	} else {
		// Use default ports
		uri.host = split1[2]
		if uri.scheme == "ws" {
			uri.port = 80
		} else {
			uri.port = 443
		}
	}

	return uri, nil
}

func parseHTTP3(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 4)
	if len(split1) < 4 {
		return nil, errors.New("incorrect URI format")
	}
	uri.scheme = split1[0][:len(split1[0])-1] // Remove trailing :
	if uri.scheme != "http" && uri.scheme != "https" {
		return nil, errors.New("unknown scheme")
	}
	uri.path = split1[3]

	// split1[1] is part between // after scheme
	split2 := strings.SplitN(split1[2], ":", 2)
	var err error
	if len(split2) == 2 {
		// Has port specified
		uri.host = split2[0]
		uri.port, err = strconv.Atoi(split2[1])
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
		if uri.port <= 0 || uri.port > 65535 {
			return nil, errors.New("port out of range")
		}
	} else {
		// Use default ports
		uri.host = split1[2]
		if uri.scheme == "http" {
			uri.port = 80
		} else {
			uri.port = 443
		}
	}

	return uri, nil
}
