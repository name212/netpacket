// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/name212/netpacket"
	stringsutils "github.com/name212/netpacket/utils/strings"
)

type Protocol uint8

// 1	Internet Control Message Protocol	ICMP
// 2	Internet Group Management Protocol	IGMP
// 6	Transmission Control Protocol	TCP
// 17	User Datagram Protocol	UDP
// 41	IPv6 encapsulation	ENCAP
// 89	Open Shortest Path First	OSPF
// 132	Stream Control Transmission Protocol	SCTP
const (
	ProtocolICMP  Protocol = 1
	ProtocolIGMP  Protocol = 2
	ProtocolTCP   Protocol = 6
	ProtocolUDP   Protocol = 17
	ProtocolENCAP Protocol = 41
	ProtocolOSPF  Protocol = 89
	ProtocolSCTP  Protocol = 132
)

var protocolsMap = map[Protocol]string{
	ProtocolICMP:  "ICMP",
	ProtocolIGMP:  "IGMP",
	ProtocolTCP:   "TCP",
	ProtocolUDP:   "UDP",
	ProtocolENCAP: "ENCAP",
	ProtocolOSPF:  "OSPF",
	ProtocolSCTP:  "SCTP",
}

// Header represents the structure of an IPv4 header
type Header struct {
	Version        uint8
	IHL            uint8
	ToS            uint8
	TotalLength    uint16
	Identification uint16
	Flags          uint8
	FragmentOffset uint16
	TTL            uint8
	Protocol       uint8
	Checksum       uint16
	SourceIP       net.IP
	DestinationIP  net.IP
	Options        []byte

	flags Flags
}

type Flags struct {
	// IsEvil
	// first bit of flags. For correct packet should be set to 0 always
	// ParsePacket check this flag before return Header
	IsEvil        bool
	DontFragment  bool
	MoreFragments bool
}

// ParseHeader parses the IPv4 header from the given byte slice
// ParseHeader save slices from data. You should copy data before parse
// to avoid hold full data in memory
func ParseHeader(data []byte) (*Header, error) {
	if err := isValidPacket(data); err != nil {
		return nil, err
	}

	header := &Header{
		Version:        data[0] >> 4,
		IHL:            extractHeaderWordsLen(data),
		ToS:            data[1],
		TotalLength:    binary.BigEndian.Uint16(data[2:4]),
		Identification: binary.BigEndian.Uint16(data[4:6]),
		Flags:          data[6] >> 5,
		FragmentOffset: binary.BigEndian.Uint16(data[6:8]) & 0x1FFF,
		TTL:            data[8],
		Protocol:       data[9],
		Checksum:       binary.BigEndian.Uint16(data[10:12]),
		SourceIP:       net.IP(data[12:16]),
		DestinationIP:  net.IP(data[16:20]),
	}

	if !header.IsValidVersion() {
		return nil, fmt.Errorf("invalid version %d", header.Version)
	}

	headerLengthBytes := header.HeaderLen()

	if header.TotalLength < minHeaderLength {
		return nil, fmt.Errorf("invalid (too small) IP total length (%d < %d)", header.TotalLength, minHeaderLength)
	}

	if headerLengthBytes < minHeaderLength {
		return nil, fmt.Errorf("invalid (too small) IP header length (%d < %d)", headerLengthBytes, minHeaderLength)
	}

	if headerLengthBytes > header.GetTotalLen() {
		return nil, fmt.Errorf("invalid IP header length > IP length (%d > %d)", headerLengthBytes, header.TotalLength)
	}

	flags := parseFlags(header.Flags)

	if flags.IsEvil {
		return nil, fmt.Errorf("invalid IP header flags first bit set to 1")
	}

	header.flags = flags

	if headerLengthBytes > 20 && len(data) >= headerLengthBytes {
		header.Options = data[20:headerLengthBytes]
	} else {
		header.Options = nil
	}

	return header, nil
}

func (h *Header) IsValidVersion() bool {
	return h.Version == 4
}

func (h *Header) GetSourceIPString() string {
	return h.SourceIP.String()
}

func (h *Header) GetSourceIP() net.IP {
	return h.SourceIP
}

func (h *Header) GetDestinationIPString() string {
	return h.DestinationIP.String()
}

func (h *Header) GetDestinationIP() net.IP {
	return h.DestinationIP
}

func (h *Header) GetTotalLen() int {
	return int(h.TotalLength)
}

func (h *Header) GetTTL() int {
	return int(h.TTL)
}

func (h *Header) GetProtocol() Protocol {
	return Protocol(h.Protocol)
}

func (h *Header) GetFlags() Flags {
	return h.flags
}

func (h *Header) ProtocolString() string {
	str, ok := protocolsMap[Protocol(h.Protocol)]
	if ok {
		return str
	}

	return "Unknown"
}

func (h *Header) HeaderLen() int {
	return headerLen(h.IHL)
}

func (h *Header) Kind() netpacket.Kind {
	return Kind
}

func (h *Header) String() string {
	s := strings.Builder{}

	flags := h.GetFlags()

	s.WriteString(stringsutils.FmtLn("Source: %s", h.SourceIP.String()))
	s.WriteString(stringsutils.FmtLn("Destination: %s", h.DestinationIP.String()))
	s.WriteString(stringsutils.FmtLn("Protocol: %s", h.ProtocolString()))
	s.WriteString(stringsutils.FmtLn("TTL: %d", h.TTL))
	s.WriteString(stringsutils.FmtLn("Header Size: %d", h.HeaderLen()))
	s.WriteString(stringsutils.FmtLn("Packet Size: %d", h.TotalLength))
	s.WriteString(stringsutils.FmtLn("Flags:"))
	s.WriteString(stringsutils.FmtLnWithTabPrefix("Don't Fragment: %v", flags.DontFragment))
	s.WriteString(stringsutils.FmtLnWithTabPrefix("More Fragments: %v", flags.MoreFragments))
	h.writeOptions(&s)
	s.WriteString(fmt.Sprintf("Checksum: %d", h.Checksum))

	return s.String()
}

func (h *Header) ParseOptions() ([]Option, error) {
	return parseOptions(h.Options)
}

func (h *Header) writeOptions(s *strings.Builder) {
	if h.Options == nil {
		s.WriteString(stringsutils.FmtLn("No options set"))
		return
	}

	opts, err := h.ParseOptions()
	if err != nil {
		s.WriteString(stringsutils.FmtLn("Cannot parse options: %v", err))
		return
	}

	s.WriteString(stringsutils.FmtLn("Options:"))

	optsStringsSlice := make([]string, 0, len(opts))

	for _, opt := range opts {
		optsStringsSlice = append(optsStringsSlice, opt.String())
	}

	optsStr := strings.Join(optsStringsSlice, "\n")

	s.WriteString(stringsutils.ShiftOnTabs(stringsutils.FmtLn(optsStr), 1))
}

func headerLen(words uint8) int {
	return int(words) * 4
}

func extractHeaderWordsLen(data []byte) uint8 {
	return data[0] & 0x0F
}

const bitSet = 1

func parseFlags(f uint8) Flags {
	evil := (f >> 0) == bitSet
	df := (f >> 1) == bitSet
	mf := (f >> 2) == bitSet

	return Flags{
		IsEvil:        evil,
		DontFragment:  df,
		MoreFragments: mf,
	}
}
