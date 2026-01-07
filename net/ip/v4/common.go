// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"fmt"

	"github.com/name212/netpacket"
)

const (
	minHeaderLength = 20

	Kind netpacket.Kind = "IPv4"
)

func isValidPacket(data []byte) error {
	if len(data) < minHeaderLength {
		return netpacket.WrapShortDataErr(fmt.Errorf("IPv4 packet"))
	}

	return nil
}
