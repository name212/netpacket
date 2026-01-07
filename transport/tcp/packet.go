// Copyright 2026
// license that can be found in the LICENSE file.

package tcp

import (
	"errors"
	"fmt"

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

func (d *Packet) GetPayload() []byte {
	return d.payload
}

func (d *Packet) GetHeader() *Header {
	return d.header
}

func (p *Packet) GetHeaderData() []byte {
	return p.headerData
}

func (d *Packet) Kind() netpacket.Kind {
	return Kind
}

func (d *Packet) GetSourcePort() int {
	// todo implement
	return 0
}

func (d *Packet) GetDestinationPort() int {
	// todo implement
	return 0
}

func (d *Packet) String() string {
	// todo implement
	return fmt.Sprintf("TCP Packet:\n\tNot implemented yet")
}

// ExtractPayload
// Warning! TODO not implemented
func ExtractPayload(data []byte) ([]byte, error) {
	return nil, netpacket.WrapNotImplementedErr(errors.New("TCP"))
}
