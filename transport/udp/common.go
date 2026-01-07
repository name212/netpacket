// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"errors"

	"github.com/name212/netpacket"
)

const (
	headerLength = 8

	Kind netpacket.Kind = "UDP"
)

func isValidDatagram(data []byte) error {
	if len(data) < headerLength {
		return netpacket.WrapShortDataErr(errors.New("UDP datagram"))
	}

	return nil
}
