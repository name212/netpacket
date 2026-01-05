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
		// Options and data would follow here if any
	}

	header, err := ParseIPv4Header(ipPacket)
	if err != nil {
		t.Error("Cannot parse header", err)
	}

	if !header.IsValidVersion() {
		t.Error("Invalid version")
	}

	if header.GetProtocol() != ProtocolTCP {
		t.Errorf("Invalid protocol: %d", header.Protocol)
	}

	if header.ProtocolString() != "TCP" {
		t.Errorf("Invalid protocol string: %s", header.ProtocolString())
	}

	if !header.SourceIP.Equal(net.ParseIP("192.168.0.104")) {
		t.Error("Invalid source ip")
	}
	if !header.DestinationIP.Equal(net.ParseIP("192.168.0.1")) {
		t.Error("Invalid source ip")
	}

	if header.TotalLength != 60 {
		t.Errorf("Invalid total length: %d", header.TotalLength)
	}

	if header.TTL != 64 {
		t.Errorf("Invalid TTL: %d", header.TTL)
	}

	flags := header.GetFlags()

	if flags.MoreFragments {
		t.Error("Invalid More fragments")
	}

	if !flags.DontFragment {
		t.Error("Invalid dont fragment")
	}

	if header.Checksum != 45542 {
		t.Errorf("Invalid checksum: %d", header.Checksum)
	}
}

func TestParseOptions(t *testing.T) {
	optionsData := []byte{
		0x01,                   // NOP
		0x07, 0x07, 0x12, 0x34, // Record Route option with example data
		0x00, // End of list
	}

	options, err := ParseOptions(optionsData)
	if err != nil {
		t.Error("cannot parse options", err)
	}

	if len(options) != 1 {
		t.Error("Invalid options")
	}

	if options[0].Type != uint8(OptionRecordRoute) {
		t.Error("Invalid option type")
	}

	if options[0].TypeShort() != "RR" {
		t.Error("Invalid option type")
	}
}
