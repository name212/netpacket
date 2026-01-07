// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/name212/netpacket"
	"github.com/name212/netpacket/utils"
)

type Header struct {
	SourcePort      uint16
	DestinationPort uint16
	Length          uint16
	Checksum        uint16
}

// ParseHeader
// header datagram from bytes
// no save any subslices from data in header
func ParseHeader(data []byte) (*Header, error) {
	if err := isValidDatagram(data); err != nil {
		return nil, err
	}

	return &Header{
		SourcePort:      binary.BigEndian.Uint16(data[0:2]),
		DestinationPort: binary.BigEndian.Uint16(data[2:4]),
		Length:          binary.BigEndian.Uint16(data[4:6]),
		Checksum:        binary.BigEndian.Uint16(data[6:8]),
	}, nil
}

func (h *Header) GetSourcePort() int {
	return int(h.SourcePort)
}

func (h *Header) GetDestinationPort() int {
	return int(h.DestinationPort)
}

func (h *Header) DatagramLen() int {
	return int(h.Length)
}

func (h *Header) HeaderLen() int {
	return headerLength
}

func (h *Header) Kind() netpacket.Kind {
	return Kind
}

func (h *Header) String() string {
	s := strings.Builder{}

	s.WriteString(utils.FmtLn("Source port: %d", h.SourcePort))
	s.WriteString(utils.FmtLn("Destination port: %d", h.DestinationPort))
	s.WriteString(utils.FmtLn("Datagram size: %d", h.Length))
	s.WriteString(fmt.Sprintf("Checksum: %d", h.Checksum))

	return s.String()
}
