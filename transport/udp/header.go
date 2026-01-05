// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"fmt"

	"github.com/name212/netpacket"
)

const minHeaderLength = 8

type Header struct {
	SourcePort uint16
	DestPort   uint16
	Length     uint16
	Checksum   uint16
}

func (h *Header) ParseHeader(data []byte) (*Header, error) {
	if len(data) < minHeaderLength {
		return nil, fmt.Errorf("%w for UDP header", ip.ShortDataErr)
	}

	return nil, nil
}
