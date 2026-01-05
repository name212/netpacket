// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"encoding/binary"
)

const headerLength = 8

type Header struct {
	SourcePort uint16
	DestPort   uint16
	Length     uint16
	Checksum   uint16
}

func (h *Header) GetSourcePort() int {
	return int(h.SourcePort)
}

func (h *Header) GetDestPort() int {
	return int(h.DestPort)
}

func (h *Header) Len() int {
	return int(h.Length)
}

func ParseHeader(data []byte) (*Header, error) {
	if err := isValidDatagram(data); err != nil {
		return nil, err
	}

	return &Header{
		SourcePort: binary.BigEndian.Uint16(data[0:2]),
		DestPort:   binary.BigEndian.Uint16(data[2:4]),
		Length:     binary.BigEndian.Uint16(data[4:6]),
		Checksum:   binary.BigEndian.Uint16(data[6:8]),
	}, nil
}
