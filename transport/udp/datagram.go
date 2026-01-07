// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"strings"

	"github.com/name212/netpacket"
	"github.com/name212/netpacket/utils"
)

type Datagram struct {
	header *Header

	headerData []byte
	payload    []byte
}

// ParseDatagram
// Parse header and extract payload from datagram
// Also save header data as subslice data
// ParseDatagram save slices from data. You should copy data before parse
// to avoid hold full data in memory
func ParseDatagram(data []byte) (*Datagram, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, netpacket.WrapCannotParseHeaderErr(err)
	}

	return &Datagram{
		header:     header,
		headerData: data[:headerLength],
		payload:    extractPayload(data),
	}, nil
}

func (d *Datagram) GetPayload() []byte {
	return d.payload
}

func (d *Datagram) GetHeader() *Header {
	return d.header
}

func (d *Datagram) GetHeaderData() []byte {
	return d.headerData
}

func (d *Datagram) Kind() netpacket.Kind {
	return Kind
}

func (d *Datagram) GetSourcePort() int {
	return d.header.GetSourcePort()
}

func (d *Datagram) GetDestinationPort() int {
	return d.header.GetDestinationPort()
}

func (d *Datagram) String() string {
	b := strings.Builder{}

	b.WriteString(utils.FmtLn("UDP Datagram:"))
	b.WriteString(utils.FmtLnWithTabPrefix("Header:"))
	b.WriteString(utils.ShiftOnTabs(d.GetHeader().String(), 2))
	b.WriteString(utils.FmtWithTabsPrefix("Payload len: %d", len(d.GetPayload())))

	return b.String()
}

// ExtractPayload extract payload from data without full parsing header
// ExtractPayload returns subslice from data. You should copy data before parse
// to avoid hold full data in memory
func ExtractPayload(data []byte) ([]byte, error) {
	if err := isValidDatagram(data); err != nil {
		return nil, err
	}

	return extractPayload(data), nil
}

func extractPayload(data []byte) []byte {
	var payload []byte
	if len(data) > headerLength {
		payload = data[headerLength:]
	}

	return payload
}
