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
}

func ParseURI(str string) (*URI, error) {
	if strings.HasPrefix(str, "udp://") || strings.HasPrefix(str, "udp4://") || strings.HasPrefix(str, "udp6://") {
		return parseUDP(str)
	} else if strings.HasPrefix(str, "tcp") || strings.HasPrefix(str, "tcp4://") || strings.HasPrefix(str, "tcp6://") {
		return parseTCP(str)
	} else if strings.HasPrefix(str, "ws") {
		return parseWebSocket(str)
	}

	return nil, errors.New("unknown scheme")
}

func parseUDP(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 3)
	if len(split1) < 3 {
		return nil, errors.New("incorrect URI format")
	}
	uri.scheme = split1[0][:len(split1[0])-1] // Exclude trailing :
	split2 := strings.SplitN(split1[2], ":", 2)
	uri.host = split2[0]
	if len(split2) < 2 {
		uri.port = 6363
	} else {
		var err error
		uri.port, err = strconv.Atoi(split2[1])
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
	}

	return uri, nil
}

func parseTCP(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 3)
	if len(split1) < 3 {
		return nil, errors.New("incorrect URI format")
	}
	uri.scheme = split1[0][:len(split1[0])-1] // Exclude trailing :
	split2 := strings.SplitN(split1[2], ":", 2)
	uri.host = split2[0]
	if len(split2) < 2 {
		uri.port = 6363
	} else {
		var err error
		uri.port, err = strconv.Atoi(split2[1])
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
	}

	return uri, nil
}

func parseWebSocket(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 3)
	if len(split1) < 3 {
		return nil, errors.New("incorrect URI format")
	}
	uri.scheme = split1[0][:len(split1[0])-1] // Exclude trailing :
	split2 := strings.SplitN(split1[2], ":", 2)
	uri.host = split2[0]
	if len(split2) < 2 {
		uri.port = 9696
	} else {
		var err error
		uri.port, err = strconv.Atoi(split2[1])
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
	}

	return uri, nil
}
