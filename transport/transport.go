/* NDN Reachability Tester
 *
 * Copyright (C) 2021 Eric Newberry.
 *
 * Released under the terms of the MIT License, as found in LICENSE.md.
 */

package transport

import "github.com/eric135/YaNFD/ndn"

type Transport interface {
	SendAndReceive(interest *ndn.Interest) error
}
