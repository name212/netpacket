// Copyright 2026
// license that can be found in the LICENSE file.

package tcp

import (
	"errors"

	"github.com/name212/netpacket"
)

type Packet struct {
	header *Header

	headerData []byte
	payload    []byte
}

// ParsePacket
// ParsePacked save s
// Warning! TODO not implemented
func ParsePacket(data []byte) (*Packet, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, err
	}
	// todo need implement
	return &Packet{
		header:  header,
		payload: nil,
	}, nil
}

func (p *Packet) GetPayload() []byte {
	return p.payload
}

func (p *Packet) GetHeader() *Header {
	return p.header
}

func (p *Packet) GetHeaderData() []byte {
	return p.headerData
}

func (p *Packet) Kind() netpacket.Kind {
	return Kind
}

func (p *Packet) GetSourcePort() int {
	// todo implement
	return 0
}

func (p *Packet) GetDestinationPort() int {
	// todo implement
	return 0
}

func (p *Packet) String() string {
	// todo implement
	return "TCP Packet:\n\tNot implemented yet"
}

// ExtractPayload extract payload from data without full parsing header
// ExtractPayload returns subslice from data. You should copy data before parse
// to avoid hold full data in memory
// Warning! TODO not implemented
func ExtractPayload(data []byte) ([]byte, error) {
	return nil, netpacket.WrapNotImplementedErr(errors.New("TCP"))
}
