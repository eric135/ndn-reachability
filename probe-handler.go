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
	"sync"
	"time"

	"github.com/eric135/YaNFD/ndn"
	"github.com/eric135/ndn-reachability/transport"
	"github.com/eric135/ndn-reachability/uri"
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
	if transportStr != "udp4" && transportStr != "udp6" && transportStr != "wss-ipv4" && transportStr != "wss-ipv6" && transportStr != "http3-ipv4" && transportStr != "http3-ipv6" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("Unknown transport"))
		return
	}
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

	wg := new(sync.WaitGroup)
	var results sync.Map
	for _, router := range routers {
		wg.Add(1)
		go p.probe(transportStr, router, interest, &results, wg)
	}
	wg.Wait()

	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(res)

	// Convert sync.Map to map
	resultsMap := make(map[string]probeResult)
	results.Range(func(key interface{}, value interface{}) bool {
		resultsMap[key.(string)] = value.(probeResult)
		return true
	})

	encoder.Encode(resultsMap)
}

func (p *ProbeHandler) probe(transportStr string, router string, interest *ndn.Interest, results *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()

	var t transport.Transport
	if transportStr == "udp4" || transportStr == "udp6" {
		uri, err := uri.ParseUDP(router)
		if err != nil {
			results.Store(router, probeResult{Ok: false, Err: "bad router"})
			return
		}
		t, err = transport.NewUDPTransport(transportStr, uri.Host, uri.Port)
		if err != nil {
			results.Store(router, probeResult{Ok: false, Err: "unable to connect"})
			return
		}
	} else if transportStr == "wss-ipv4" || transportStr == "wss-ipv6" {
		var err error
		t, err = transport.NewWebSocketTransport(transportStr, router)
		if err != nil {
			results.Store(router, probeResult{Ok: false, Err: "unable to connect"})
			return
		}
	} else if transportStr == "http3-ipv4" || transportStr == "http3-ipv6" {
		// TODO
	}

	// Send and receive
	if rtt, err := t.SendAndReceive(interest); err != nil {
		results.Store(router, probeResult{Ok: false, Err: err.Error()})
	} else {
		results.Store(router, probeResult{Ok: true, RTT: uint(rtt.Milliseconds())})
	}
}
