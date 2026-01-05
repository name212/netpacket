// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import "github.com/name212/netpacket"

func isValidDatagram(data []byte) error {
	if len(data) < headerLength {
		return netpacket.WrapShortDataErr("UDP datagram")
	}

	return nil
}
