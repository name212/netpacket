// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"fmt"

	"github.com/name212/netpacket"
)

type Packet struct {
	Header  *Header
	Payload []byte
}

func ParsePacket(data []byte) (*Packet, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, netpacket.WrapCannotParseHeaderErr(err.Error())
	}

	totalLen := header.GetTotalLen()

	if totalLen > len(data) {
		return nil, fmt.Errorf("data too short to contain an IPv4 all packet header len %d data len %d", totalLen, len(data))
	}

	var payload []byte
	headerLengthBytes := header.HeaderLen()
	if totalLen > headerLengthBytes {
		payload = data[headerLengthBytes:]
	}

	return &Packet{
		Header:  header,
		Payload: payload,
	}, nil
}
