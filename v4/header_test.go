// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"net"
	"testing"
)

func TestParseHeader(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x3c,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	header, err := ParseIPv4Header(ipPacket)
	if err != nil {
		t.Error("Cannot parse header", err)
		return
	}

	if !header.IsValidVersion() {
		t.Error("Invalid version")
		return
	}

	if header.GetProtocol() != ProtocolTCP {
		t.Errorf("Invalid protocol: %d", header.Protocol)
		return
	}

	if header.ProtocolString() != "TCP" {
		t.Errorf("Invalid protocol string: %s", header.ProtocolString())
		return
	}

	if !header.SourceIP.Equal(net.ParseIP("192.168.0.104")) {
		t.Error("Invalid source ip")
		return
	}

	if !header.DestinationIP.Equal(net.ParseIP("192.168.0.1")) {
		t.Error("Invalid source ip")
		return
	}

	if header.TotalLength != 60 {
		t.Errorf("Invalid total length: %d", header.TotalLength)
		return
	}

	if header.TTL != 64 {
		t.Errorf("Invalid TTL: %d", header.TTL)
		return
	}

	flags := header.GetFlags()

	if flags.IsEvil {
		t.Error("Evil set")
		return
	}

	if flags.MoreFragments {
		t.Error("Invalid More fragments")
		return
	}

	if !flags.DontFragment {
		t.Error("Invalid dont fragment")
		return
	}

	if header.Checksum != 45542 {
		t.Errorf("Invalid checksum: %d", header.Checksum)
		return
	}
}

func TestParseOptions(t *testing.T) {
	ipPacket := []byte{
		0x49, 0x00, 0x00, 0x28, 0x03, 0x04,
		0x00, 0x00, 0xfe, 0x01, 0xc1, 0xe0,
		0xaf, 0x2d, 0xb0, 0x00, 0x95, 0xab,
		0x7e, 0x0b, 0x01, 0x82, 0x0b, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
	}

	header, err := ParseIPv4Header(ipPacket)
	if err != nil {
		t.Error("Cannot parse header", err)
		return
	}

	if header.ProtocolString() != "ICMP" {
		t.Errorf("Invalid protocol string: %s", header.ProtocolString())
		return
	}

	options, err := header.ParseOptions()
	if err != nil {
		t.Error("Cannot parse options", err)
		return
	}

	if len(options) != 2 {
		t.Error("Invalid options len", options)
	}

	firstOption := options[0]

	if firstOption.GetType() != OptionNoOperation {
		t.Error("Invalid first option. Should NoOp", firstOption.GetType())
		return
	}

	if firstOption.Length != 1 {
		t.Error("Invalid first option len. Should 1", firstOption.Length)
		return
	}

	if len(firstOption.Data) != 0 {
		t.Error("Invalid first option NoOp should not have data", len(firstOption.Data))
		return
	}

	if firstOption.TypeShort() != "NOP" {
		t.Error("Invalid first option. Should NOP", firstOption.TypeShort())
	}

	secondOption := options[1]

	if secondOption.GetType() != OptionSecurityRIPSO {
		t.Error("Invalid second option. Should Security RIPSO", secondOption.GetType())
		return
	}

	if len(secondOption.Data) != 9 {
		t.Error("Invalid second option. Should have len 9", len(firstOption.Data))
		return
	}

	if secondOption.TypeShort() != "SEC" {
		t.Error("Invalid second option. Should SEC", secondOption.TypeShort())
	}
}
