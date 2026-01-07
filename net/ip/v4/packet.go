// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/name212/netpacket"
	"github.com/name212/netpacket/transport/tcp"
	"github.com/name212/netpacket/transport/udp"
	stringsutils "github.com/name212/netpacket/utils/strings"
)

var ErrNotTransportPacket = errors.New("not transport packet")

type Transport interface {
	GetSourcePort() int
	GetDestinationPort() int
	GetPayload() []byte
	Kind() netpacket.Kind
}

type Packet struct {
	header *Header

	headerData []byte
	payload    []byte
}

// ParsePacket parses the IPv4 header and extract payload also save header data
// ParsePacket save slices from data. You should copy data before parse
// to avoid hold full data in memory
func ParsePacket(data []byte) (*Packet, error) {
	header, err := ParseHeader(data)
	if err != nil {
		return nil, netpacket.WrapCannotParseHeaderErr(err)
	}

	totalLen := header.GetTotalLen()

	if totalLen > len(data) {
		return nil, fmt.Errorf("data too short to contain an IPv4 all packet header len %d data len %d", totalLen, len(data))
	}

	return &Packet{
		header:     header,
		headerData: data[:header.HeaderLen()],
		payload:    extractPayload(data, header.HeaderLen()),
	}, nil
}

func (p *Packet) GetSourceIP() net.IP {
	return p.GetHeader().GetSourceIP()
}

func (p *Packet) GetSourceIPString() string {
	return p.GetHeader().GetSourceIPString()
}

func (p *Packet) GetDestinationIPString() string {
	return p.GetHeader().GetDestinationIPString()
}

func (p *Packet) GetDestinationIP() net.IP {
	return p.GetHeader().GetDestinationIP()
}

func (p *Packet) GetTTL() int {
	return p.GetHeader().GetTTL()
}

func (p *Packet) GetProtocol() Protocol {
	return p.GetHeader().GetProtocol()
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

func (p *Packet) IsTransport() bool {
	proto := p.GetHeader().GetProtocol()
	return proto == ProtocolTCP || proto == ProtocolUDP
}

// TransportPacket
// returns ErrNotTransportPacket error if packet is not UDP or TCP
func (p *Packet) TransportPacket() (Transport, error) {
	payload := p.GetPayload()
	if len(payload) == 0 {
		return nil, netpacket.WrapShortDataErr(netpacket.ErrEmptyPayload)
	}

	header := p.GetHeader()

	var inner Transport
	var err error

	switch header.GetProtocol() {
	case ProtocolTCP:
		inner, err = tcp.ParsePacket(payload)
	case ProtocolUDP:
		inner, err = udp.ParseDatagram(payload)
	default:
		return nil, fmt.Errorf("%w %s", ErrNotTransportPacket, header.ProtocolString())
	}

	if err != nil {
		return nil, err
	}

	return inner, nil
}

func (p *Packet) String() string {
	b := strings.Builder{}

	b.WriteString(stringsutils.FmtLn("IPv4 Packet:"))
	b.WriteString(stringsutils.FmtLnWithTabPrefix("Header:"))
	b.WriteString(stringsutils.ShiftOnTabs(stringsutils.FmtLn(p.GetHeader().String()), 2))
	b.WriteString(stringsutils.FmtLnWithTabPrefix("Is transport: %v", p.IsTransport()))
	b.WriteString(stringsutils.FmtWithTabPrefix("Payload len: %d", len(p.GetPayload())))

	return b.String()
}

// ToUDP
// Warning! no additional checks before convert. Can panic.
// Please check Transport.Kind before conversion
func ToUDP(t Transport) *udp.Datagram {
	return t.(*udp.Datagram)
}

// ToTCP
// Warning! no additional checks before convert. Can panic.
// Please check Transport.Kind before conversion
func ToTCP(t Transport) *tcp.Packet {
	return t.(*tcp.Packet)
}

// ExtractPayload extract payload from data without full parsing header
// ExtractPayload returns subslice from data. You should copy data before parse
// to avoid hold full data in memory
func ExtractPayload(data []byte) ([]byte, error) {
	if err := isValidPacket(data); err != nil {
		return nil, err
	}

	headerLenWords := extractHeaderWordsLen(data)

	return extractPayload(data, headerLen(headerLenWords)), nil
}

func extractPayload(data []byte, headerLengthBytes int) []byte {
	var payload []byte
	if len(data) > headerLengthBytes {
		payload = data[headerLengthBytes:]
	}

	return payload
}
