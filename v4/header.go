// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"github.com/name212/ip/utils"
)

const (
	minHeaderLength = 20
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
	IsEvil        bool
	DontFragment  bool
	MoreFragments bool
}

// ParseHeader parses the IPv4 header from the given byte slice
func ParseHeader(data []byte) (*Header, error) {
	if len(data) < minHeaderLength {
		return nil, fmt.Errorf("data too short to contain an IPv4 header")
	}

	header := &Header{
		Version:        data[0] >> 4,
		IHL:            data[0] & 0x0F,
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

func (h *Header) GetTotalLen() int {
	return int(h.TotalLength)
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
	return int(h.IHL) * 4
}

func (h *Header) String() string {
	s := strings.Builder{}

	flags := h.GetFlags()

	s.WriteString(utils.FmtLn("Source: %s", h.SourceIP.String()))
	s.WriteString(utils.FmtLn("Destination: %s", h.DestinationIP.String()))
	s.WriteString(utils.FmtLn("Protocol: %s", h.ProtocolString()))
	s.WriteString(utils.FmtLn("TTL: %d", h.TTL))
	s.WriteString(utils.FmtLn("Size: %d", h.TotalLength))
	s.WriteString(utils.FmtLn("Flags:"))
	s.WriteString(utils.FmtLnWithTabPrefix("Don't Fragment: %v", flags.DontFragment))
	s.WriteString(utils.FmtLnWithTabPrefix("More Fragments: %v", flags.MoreFragments))
	h.writeOptions(&s)
	s.WriteString(utils.FmtLn("Checksum: %d", h.Checksum))

	return s.String()
}

func (h *Header) ParseOptions() ([]Option, error) {
	if h.Options == nil {
		return nil, nil
	}

	data := h.Options
	res := make([]Option, 0, 4)

	for len(data) > 0 {
		opt := Option{Type: data[0]}

		switch opt.GetType() {
		case OptionEndOfList: // End of options
			return res, nil
		case OptionNoOperation: // 1 byte padding
			opt.Length = 1
			data = data[1:]
			res = append(res, opt)
		default:
			shortWithIDDesc := opt.TypeShortWithID()
			if len(data) < 2 {
				return nil, fmt.Errorf("invalid option length. Length %d less than 2 for %s", len(data), shortWithIDDesc)
			}
			opt.Length = data[1]
			if len(data) < int(opt.Length) {
				return nil, fmt.Errorf("IP option length exceeds remaining IP header size, option for %s length %v", shortWithIDDesc, opt.Length)
			}
			if opt.Length <= 2 {
				return nil, fmt.Errorf("invalid IP option type %s length %d. Must be greater than 2", shortWithIDDesc, opt.Length)
			}
			opt.Data = data[2:opt.Length]
			data = data[opt.Length:]
			res = append(res, opt)
		}
	}

	return res, nil
}

func (h *Header) writeOptions(s *strings.Builder) {
	if h.Options == nil {
		s.WriteString(utils.FmtLn("No options set"))
		return
	}

	opts, err := h.ParseOptions()
	if err != nil {
		s.WriteString(utils.FmtLn("Cannot parse options: %v", err))
		return
	}

	s.WriteString(utils.FmtLn("Options:"))
	for _, opt := range opts {
		msg := utils.FmtLnWithTabPrefix("%s: data len %d", opt.TypeLong(), len(opt.Data))
		s.WriteString(msg)
	}
}

func parseFlags(f uint8) Flags {
	evil := utils.CheckBit(f, 0)
	df := utils.CheckBit(f, 1)
	mf := utils.CheckBit(f, 2)

	return Flags{
		IsEvil:        evil,
		DontFragment:  df,
		MoreFragments: mf,
	}
}
