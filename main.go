/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := "80"
	if len(os.Args) >= 2 {
		port = os.Args[1]
		portInt, err := strconv.Atoi(os.Args[1])
		if err != nil || portInt <= 0 || portInt > 65535 {
			fmt.Println("Usage:")
			fmt.Println(os.Args[0], "[port]")
			os.Exit(-1)
		}
	}

	probeHandler := new(ProbeHandler)

	http.Handle("/probe", probeHandler)
	http.ListenAndServe(":"+port, nil)
}
