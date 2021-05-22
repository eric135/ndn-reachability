/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eric135/YaNFD/ndn"
	"github.com/eric135/ndn-reachability/transport"
)

type ProbeHandler struct {
}

type probeResult struct {
	Ok  bool   `json:"ok"`
	RTT uint   `json:"rtt,omitempty"`
	Err string `json:"error,omitempty"`
}

func (p *ProbeHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		// We only want to process POST
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unsupported method"))
		return
	}

	req.ParseForm()

	if req.Form.Get("transport") == "" || req.Form.Get("router") == "" || req.Form.Get("name") == "" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Missing field"))
		return
	}

	transportStr := req.Form.Get("transport")
	routers := req.Form["router"]
	nameStr := req.Form.Get("name")
	suffix := req.Form.Get("suffix")
	if suffix != "1" {
		suffix = "0"
	}
	name, err := ndn.NameFromString(nameStr)
	if err != nil {
		// Incorrectly-formatted name
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Incorrect name format"))
		return
	}
	if suffix != "0" {
		name.Append(ndn.NewTimestampNameComponent(uint64(time.Now().Nanosecond())))
	}

	interest := ndn.NewInterest(name)

	results := make(map[string]probeResult)
	for _, router := range routers {
		var t transport.Transport
		if transportStr == "udp4" || transportStr == "udp6" {
			uri, err := parseUDP(router)
			if err != nil {
				res.WriteHeader(http.StatusBadRequest)
				res.Write([]byte("Bad URI"))
				return
			}
			t = transport.NewUDPTransport(transportStr, uri.host, uri.port)
		} else if transportStr == "wss-ipv4" || transportStr == "wss-ipv6" {
			// TODO
		} else if transportStr == "http3-ipv4" || transportStr == "http3-ipv6" {
			// TODO
		} else {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Unknown transport"))
			return
		}

		// If nil, error establishing connection
		if t == nil {
			results[router] = probeResult{Ok: false, Err: "internal error"}
		}

		// Send and receive
		if rtt, err := t.SendAndReceive(interest); err != nil {
			results[router] = probeResult{Ok: false, Err: err.Error()}
		} else {
			results[router] = probeResult{Ok: true, RTT: uint(rtt.Milliseconds())}
		}
	}

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(res)
	encoder.Encode(results)
}
