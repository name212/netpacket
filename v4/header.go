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
}

type Flags struct {
	DontFragment  bool
	MoreFragments bool
}

func (h *Header) IsValidVersion() bool {
	return h.Version == 4
}

func (h *Header) GetProtocol() Protocol {
	return Protocol(h.Protocol)
}

func (h *Header) GetFlags() Flags {
	df := utils.CheckBit(h.Flags, 1)
	mf := utils.CheckBit(h.Flags, 2)

	return Flags{
		DontFragment:  df,
		MoreFragments: mf,
	}
}

func (h *Header) ProtocolString() string {
	str, ok := protocolsMap[Protocol(h.Protocol)]
	if ok {
		return str
	}

	return "Unknown"
}

func (h *Header) writeOptions(s *strings.Builder) {
	if h.Options == nil {
		s.WriteString(utils.FmtLn("No options set"))
		return
	}

	opts, err := ParseOptions(h.Options)
	if err != nil {
		s.WriteString(utils.FmtLn("Cannot parse options: %v", err))
		return
	}

	for _, opt := range opts {
		msg := utils.FmtLnWithTabPrefix("%s: data len %d", opt.TypeLong(), len(opt.Data))
		s.WriteString(msg)
	}
}

func (h *Header) String() string {
	s := strings.Builder{}

	flags := h.GetFlags()

	s.WriteString(utils.FmtLn("IsValidVersion: %v", h.IsValidVersion()))
	s.WriteString(utils.FmtLn("Protocol: %s", h.ProtocolString()))
	s.WriteString(utils.FmtLn("Source: %s", h.SourceIP.String()))
	s.WriteString(utils.FmtLn("Destination: %s", h.DestinationIP.String()))
	s.WriteString(utils.FmtLn("TTL: %d", h.TTL))
	s.WriteString(utils.FmtLn("Size: %d", h.TotalLength))
	s.WriteString(utils.FmtLn("Flags:"))
	s.WriteString(utils.FmtLnWithTabPrefix("Don't Fragment: %v", flags.DontFragment))
	s.WriteString(utils.FmtLnWithTabPrefix("More Fragments: %v", flags.MoreFragments))
	h.writeOptions(&s)
	s.WriteString(utils.FmtLn("Checksum: %d", h.Checksum))

	return s.String()
}

// ParseIPv4Header parses the IPv4 header from the given byte slice
func ParseIPv4Header(data []byte) (*Header, error) {
	if len(data) < 20 {
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

	// Calculate the header length in bytes
	headerLength := int(header.IHL) * 4

	// Check if options exist
	if headerLength > 20 && len(data) >= headerLength {
		header.Options = data[20:headerLength]
	} else {
		header.Options = nil
	}

	return header, nil
}

type OptionType uint8

// 0/0x00	EOOL	End of Option List
// 1/0x01	NOP	No Operation
// 2/0x02	SEC	Security (defunct)
// 7/0x07	RR	Record Route
// 10/0x0A	ZSU	Experimental Measurement
// 11/0x0B	MTUP	MTU Probe
// 12/0x0C	MTUR	MTU Reply
// 15/0x0F	ENCODE	ENCODE
// 25/0x19	QS	Quick-Start
// 30/0x1E	EXP	RFC3692-style Experiment
// 68/0x44	TS	Time Stamp
// 82/0x52	TR	Traceroute
// 94/0x5E	EXP	RFC3692-style Experiment
// 130/0x82	SEC	Security (RIPSO)
// 131/0x83	LSR	Loose Source Route
// 133/0x85	E-SEC	Extended Security (RIPSO)
// 134/0x86	CIPSO	Commercial IP Security Option
// 136/0x88	SID	Stream ID
// 137/0x89	SSR	Strict Source Route
// 142/0x8E	VISA	Experimental Access Control
// 144/0x90	IMITD	IMI Traffic Descriptor
// 145/0x91	EIP	Extended Internet Protocol
// 147/0x93	ADDEXT	Address Extension
// 148/0x94	RTRALT	Router Alert
// 149/0x95	SDB	Selective Directed Broadcast
// 151/0x97	DPS	Dynamic Packet State
// 152/0x98	UMP	Upstream Multicast Packet
// 158/0x9E	EXP	RFC3692-style Experiment
// 205/0xCD	FINN	Experimental Flow Control
// 222/0xDE	EXP	RFC3692-style Experiment
const (
	OptionEndOfList                    OptionType = 0
	OptionNoOperation                  OptionType = 1
	OptionSecurityDefunct              OptionType = 2
	OptionRecordRoute                  OptionType = 7
	OptionExperimentalMeasurement      OptionType = 10
	OptionMTUProbe                     OptionType = 11
	OptionMTUReply                     OptionType = 12
	OptionENCODE                       OptionType = 15
	OptionQuickStart                   OptionType = 25
	OptionRFC3692StyleExperimentFirst  OptionType = 30
	OptionTimeStamp                    OptionType = 68
	OptionTraceroute                   OptionType = 82
	OptionRFC3692StyleExperimentSecond OptionType = 94
	OptionSecurityRIPSO                OptionType = 130
	OptionLooseSourceRoute             OptionType = 131
	OptionExtendedSecurityRIPSO        OptionType = 133
	OptionCommercialIPSecurityOption   OptionType = 134
	OptionStreamID                     OptionType = 136
	OptionStrictSourceRoute            OptionType = 137
	OptionExperimentalAccessControl    OptionType = 142
	OptionIMITrafficDescriptor         OptionType = 144
	OptionExtendedInternetProtocol     OptionType = 145
	OptionAddressExtension             OptionType = 147
	OptionRouterAlert                  OptionType = 148
	OptionSelectiveDirectedBroadcast   OptionType = 149
	OptionDynamicPacketState           OptionType = 151
	OptionUpstreamMulticastPacket      OptionType = 152
	OptionRFC3692StyleExperimentThird  OptionType = 158
	OptionExperimentalFlowControl      OptionType = 205
	OptionRFC3692StyleExperimentFour   OptionType = 222
)

type optionDescription struct {
	short string
	long  string
}

var optionTypesMap = map[OptionType]*optionDescription{
	OptionEndOfList: &optionDescription{
		short: "EOOL",
		long:  "End Of List",
	},
	OptionNoOperation: &optionDescription{
		short: "NOP",
		long:  "No Operation",
	},
	OptionSecurityDefunct: &optionDescription{
		short: "SEC",
		long:  "Security Defunct",
	},
	OptionRecordRoute: &optionDescription{
		short: "ROR",
		long:  "Record Route",
	},
	OptionExperimentalMeasurement: &optionDescription{
		short: "EXP",
		long:  "Experimental Measurement",
	},
	OptionMTUProbe: &optionDescription{
		short: "MTUP",
		long:  "MTU Probe",
	},
	OptionMTUReply: &optionDescription{
		short: "MTUR",
		long:  "MTU Reply",
	},
	OptionENCODE: &optionDescription{
		short: "ENCODE",
		long:  "ENCODE",
	},
	OptionQuickStart: &optionDescription{
		short: "QS",
		long:  "Quick Start",
	},
	OptionRFC3692StyleExperimentFirst: &optionDescription{
		short: "EXP",
		long:  "RFC3692 Experiment First",
	},
	OptionTimeStamp: &optionDescription{
		short: "TS",
		long:  "Timestamp",
	},
	OptionTraceroute: &optionDescription{
		short: "TR",
		long:  "Traceroute",
	},
	OptionRFC3692StyleExperimentSecond: &optionDescription{
		short: "EXP",
		long:  "RFC3692 Experiment Second",
	},
	OptionSecurityRIPSO: &optionDescription{
		short: "SEC",
		long:  "Security RIPSO",
	},
	OptionLooseSourceRoute: &optionDescription{
		short: "LSR",
		long:  "Loose Source Route",
	},
	OptionExtendedSecurityRIPSO: &optionDescription{
		short: "E-SEC",
		long:  "Extended Security (RIPSO)",
	},
	OptionCommercialIPSecurityOption: &optionDescription{
		short: "CIPSO",
		long:  "Commercial IP Security Option",
	},
	OptionStreamID: &optionDescription{
		short: "SID",
		long:  "Stream ID",
	},
	OptionStrictSourceRoute: &optionDescription{
		short: "SSR",
		long:  "Strict Source Route",
	},
	OptionExperimentalAccessControl: &optionDescription{
		short: "VISA",
		long:  "Experimental Access Control",
	},
	OptionIMITrafficDescriptor: &optionDescription{
		short: "IMITD",
		long:  "IMI Traffic Descriptor",
	},
	OptionExtendedInternetProtocol: &optionDescription{
		short: "EIP",
		long:  "Extended Internet Protocol",
	},
	OptionAddressExtension: &optionDescription{
		short: "ADDEXT",
		long:  "Address Extension",
	},
	OptionRouterAlert: &optionDescription{
		short: "RTRALT",
		long:  "Router Alert",
	},
	OptionSelectiveDirectedBroadcast: &optionDescription{
		short: "SDB",
		long:  "Selective Directed Broadcast",
	},
	OptionDynamicPacketState: &optionDescription{
		short: "DPS",
		long:  "Dynamic Packet State",
	},
	OptionUpstreamMulticastPacket: &optionDescription{
		short: "UMP",
		long:  "Upstream Multicast Packet",
	},
	OptionRFC3692StyleExperimentThird: &optionDescription{
		short: "EXP",
		long:  "RFC3692 Experiment Third",
	},
	OptionExperimentalFlowControl: &optionDescription{
		short: "EXPF",
		long:  "Experimental Flow Control",
	},
	OptionRFC3692StyleExperimentFour: &optionDescription{
		short: "EXP",
		long:  "RFC3692 Experiment Four",
	},
}

func unknownOptionDescription(optionType OptionType) *optionDescription {
	return &optionDescription{
		short: "UNKNOWN",
		long:  fmt.Sprintf("Unknown: %d", optionType),
	}
}

type Option struct {
	Type   uint8
	Length uint8
	Data   []byte
}

func (o *Option) TypeShort() string {
	desc, ok := optionTypesMap[OptionType(o.Type)]
	if ok {
		return desc.short
	}

	return unknownOptionDescription(OptionType(o.Type)).short
}

func (o *Option) TypeLong() string {
	desc, ok := optionTypesMap[OptionType(o.Type)]
	if ok {
		return desc.long
	}

	return unknownOptionDescription(OptionType(o.Type)).long
}

func (o *Option) String() string {
	return fmt.Sprintf("Option Type: %s(%d), Length: %d, Data: %v", o.TypeShort(), o.Type, o.Length, o.Data)
}

func ParseOptions(data []byte) ([]Option, error) {
	var options []Option
	offset := 0

	for offset < len(data) {
		if len(data)-offset < 1 {
			return nil, fmt.Errorf("insufficient data for parsing options")
		}

		optionType := data[offset]
		offset++

		if optionType == 0 { // End of options list
			break
		} else if optionType == 1 { // No operation (NOP)
			options = append(options, Option{Type: optionType})
			continue
		} else {
			if len(data)-offset < 1 {
				return nil, fmt.Errorf("insufficient data for option length")
			}
			optionLength := data[offset]
			offset++

			if optionLength < 2 || offset+int(optionLength)-2 > len(data) {
				return nil, fmt.Errorf("invalid option length")
			}

			optionData := data[offset : offset+int(optionLength)-2]
			offset += int(optionLength) - 2

			options = append(options, Option{
				Type:   optionType,
				Length: optionLength,
				Data:   optionData,
			})
		}
	}

	return options, nil
}
