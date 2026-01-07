// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"fmt"
	"strings"

	stringsutils "github.com/name212/netpacket/utils/strings"
)

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

type Option struct {
	typeID uint8
	length uint8
	data   []byte
}

func parseOptions(data []byte) ([]Option, error) {
	if data == nil {
		return nil, nil
	}

	res := make([]Option, 0, 4)

	for len(data) > 0 {
		opt := Option{typeID: data[0]}

		switch opt.GetType() {
		case OptionEndOfList:
			return res, nil
		case OptionNoOperation:
			opt.length = 1
			data = data[1:]
			res = append(res, opt)
		default:
			if len(data) < 2 {
				return nil, opt.wrapError("invalid length. Length %d less than 2", len(data))
			}
			opt.length = data[1]
			intLen := opt.GetLength()
			if len(data) < intLen {
				return nil, opt.wrapError("length exceeds remaining IP header size, length %v", intLen)
			}
			if intLen <= 2 {
				return nil, opt.wrapError("invalid length %d. Must be greater than 2", intLen)
			}
			opt.data = data[2:intLen]
			data = data[intLen:]
			res = append(res, opt)
		}
	}

	return res, nil
}

func (o *Option) GetLength() int {
	return int(o.length)
}

func (o *Option) GetData() []byte {
	return o.data
}

func (o *Option) GetType() OptionType {
	return OptionType(o.typeID)
}

func (o *Option) TypeShort() string {
	return getOptionDescription(o.GetType()).short
}

func (o *Option) TypeShortWithID() string {
	return fmt.Sprintf("%s(%d)", o.TypeShort(), o.GetType())
}

func (o *Option) TypeLong() string {
	return getOptionDescription(o.GetType()).long
}

func (o *Option) writeData(b *strings.Builder) {
	data := o.GetData()
	if len(data) == 0 {
		b.WriteString("\tNo data")
		return
	}

	b.WriteString(stringsutils.FmtLnWithTabPrefix("Hex data:"))
	b.WriteString(stringsutils.ShiftOnTabs(stringsutils.BytesToHexWithWrap(data, 8), 2))
}

func (o *Option) String() string {
	b := strings.Builder{}

	b.WriteString(stringsutils.FmtLn("Option:"))
	b.WriteString(stringsutils.FmtLnWithTabPrefix("Type: %s", o.TypeShortWithID()))
	b.WriteString(stringsutils.FmtLnWithTabPrefix("Type description: %s", o.TypeLong()))
	b.WriteString(stringsutils.FmtLnWithTabPrefix("Full Length: %d", o.GetLength()))
	o.writeData(&b)

	return b.String()
}

func (o *Option) wrapError(f string, args ...any) error {
	f = fmt.Sprintf("option %s: ", o.TypeShortWithID()) + f
	return fmt.Errorf(f, args...)
}

type optionDescription struct {
	short string
	long  string
}

var optionTypesMap = map[OptionType]*optionDescription{
	OptionEndOfList: {
		short: "EOOL",
		long:  "End Of List",
	},
	OptionNoOperation: {
		short: "NOP",
		long:  "No Operation",
	},
	OptionSecurityDefunct: {
		short: "SEC",
		long:  "Security Defunct",
	},
	OptionRecordRoute: {
		short: "ROR",
		long:  "Record Route",
	},
	OptionExperimentalMeasurement: {
		short: "EXP",
		long:  "Experimental Measurement",
	},
	OptionMTUProbe: {
		short: "MTUP",
		long:  "MTU Probe",
	},
	OptionMTUReply: {
		short: "MTUR",
		long:  "MTU Reply",
	},
	OptionENCODE: {
		short: "ENCODE",
		long:  "ENCODE",
	},
	OptionQuickStart: {
		short: "QS",
		long:  "Quick Start",
	},
	OptionRFC3692StyleExperimentFirst: {
		short: "EXP",
		long:  "RFC3692 Experiment First",
	},
	OptionTimeStamp: {
		short: "TS",
		long:  "Timestamp",
	},
	OptionTraceroute: {
		short: "TR",
		long:  "Traceroute",
	},
	OptionRFC3692StyleExperimentSecond: {
		short: "EXP",
		long:  "RFC3692 Experiment Second",
	},
	OptionSecurityRIPSO: {
		short: "SEC",
		long:  "Security RIPSO",
	},
	OptionLooseSourceRoute: {
		short: "LSR",
		long:  "Loose Source Route",
	},
	OptionExtendedSecurityRIPSO: {
		short: "E-SEC",
		long:  "Extended Security (RIPSO)",
	},
	OptionCommercialIPSecurityOption: {
		short: "CIPSO",
		long:  "Commercial IP Security Option",
	},
	OptionStreamID: {
		short: "SID",
		long:  "Stream ID",
	},
	OptionStrictSourceRoute: {
		short: "SSR",
		long:  "Strict Source Route",
	},
	OptionExperimentalAccessControl: {
		short: "VISA",
		long:  "Experimental Access Control",
	},
	OptionIMITrafficDescriptor: {
		short: "IMITD",
		long:  "IMI Traffic Descriptor",
	},
	OptionExtendedInternetProtocol: {
		short: "EIP",
		long:  "Extended Internet Protocol",
	},
	OptionAddressExtension: {
		short: "ADDEXT",
		long:  "Address Extension",
	},
	OptionRouterAlert: {
		short: "RTRALT",
		long:  "Router Alert",
	},
	OptionSelectiveDirectedBroadcast: {
		short: "SDB",
		long:  "Selective Directed Broadcast",
	},
	OptionDynamicPacketState: {
		short: "DPS",
		long:  "Dynamic Packet State",
	},
	OptionUpstreamMulticastPacket: {
		short: "UMP",
		long:  "Upstream Multicast Packet",
	},
	OptionRFC3692StyleExperimentThird: {
		short: "EXP",
		long:  "RFC3692 Experiment Third",
	},
	OptionExperimentalFlowControl: {
		short: "EXPF",
		long:  "Experimental Flow Control",
	},
	OptionRFC3692StyleExperimentFour: {
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

func getOptionDescription(optionType OptionType) *optionDescription {
	description, ok := optionTypesMap[optionType]
	if ok {
		return description
	}

	return unknownOptionDescription(optionType)
}
