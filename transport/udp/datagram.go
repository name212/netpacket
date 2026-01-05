// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import "github.com/name212/netpacket"

type Datagram struct {
	Header  *Header
	Payload []byte
}

func ParseDatagram(data []byte) (*Datagram, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, netpacket.WrapCannotParseHeaderErr(err.Error())
	}

	return &Datagram{
		Header:  header,
		Payload: extractPayload(data),
	}, nil
}

func extractPayload(data []byte) []byte {
	var payload []byte
	if len(data) > headerLength {
		payload = data[headerLength:]
	}

	return payload
}

func ExtractPayload(data []byte) ([]byte, error) {
	if err := isValidDatagram(data); err != nil {
		return nil, err
	}

	return extractPayload(data), nil
}
