/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package uri

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

type URI struct {
	Scheme string
	Host   string
	Port   int
	Path   string
}

func ParseUDP(str string) (*URI, error) {
	uri := new(URI)
	uri.Scheme = "udp"
	split := strings.SplitN(str, ":", 2)
	if len(split) < 2 {
		return nil, errors.New("incorrect URI format")
	}
	uri.Host = split[0]
	var err error
	uri.Port, err = strconv.Atoi(split[1])
	if err != nil {
		return nil, errors.New("could not convert port to integer")
	}
	if uri.Port <= 0 || uri.Port > 65535 {
		return nil, errors.New("port out of range")
	}
	return uri, nil
}

func ParseWebSocket(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 4)
	if len(split1) < 4 {
		return nil, errors.New("incorrect URI format")
	}
	uri.Scheme = split1[0][:len(split1[0])-1] // Remove trailing :
	if uri.Scheme != "ws" && uri.Scheme != "wss" {
		return nil, errors.New("unknown scheme")
	}
	uri.Path = split1[3]

	// split1[1] is part between // after scheme
	var err error
	var portStr string
	uri.Host, portStr, err = net.SplitHostPort(split1[2])
	if err != nil {
		// Has no port specified - use defaults
		uri.Host = split1[2]
		uri.Port = 9696
	} else {
		uri.Port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
		if uri.Port <= 0 || uri.Port > 65535 {
			return nil, errors.New("port out of range")
		}
	}

	return uri, nil
}

func ParseHTTP3(str string) (*URI, error) {
	uri := new(URI)
	split1 := strings.SplitN(str, "/", 4)
	if len(split1) < 4 {
		return nil, errors.New("incorrect URI format")
	}
	uri.Scheme = split1[0][:len(split1[0])-1] // Remove trailing :
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return nil, errors.New("unknown scheme")
	}
	uri.Path = split1[3]

	// split1[1] is part between // after scheme
	var err error
	var portStr string
	uri.Host, portStr, err = net.SplitHostPort(split1[2])
	if err != nil {
		// Has no port specified - use defaults
		uri.Host = split1[2]
		if uri.Scheme == "http" {
			uri.Port = 80
		} else {
			uri.Port = 443
		}
	} else {
		uri.Port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, errors.New("could not convert port to integer")
		}
		if uri.Port <= 0 || uri.Port > 65535 {
			return nil, errors.New("port out of range")
		}
	}

	return uri, nil
}
