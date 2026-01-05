// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import "fmt"

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
	Type   uint8
	Length uint8
	Data   []byte
}

func (o *Option) GetType() OptionType {
	return OptionType(o.Type)
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

func (o *Option) String() string {
	return fmt.Sprintf("Option Type: %s(%d), Length: %d, Data: %v", o.TypeShort(), o.Type, o.Length, o.Data)
}

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

func getOptionDescription(optionType OptionType) *optionDescription {
	description, ok := optionTypesMap[optionType]
	if ok {
		return description
	}

	return unknownOptionDescription(optionType)
}
