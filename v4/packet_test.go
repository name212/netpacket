// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"net"
	"testing"
)

func TestParsePacketFailIfDataShort(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x3c,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	_, err := ParsePacket(ipPacket)
	if err == nil {
		t.Error("Packet should not parsed", err)
		return
	}
}

func TestParsePacketNoData(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x14,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	packet, err := ParsePacket(ipPacket)
	if err != nil {
		t.Error("Packet should parsed", err)
		return
	}

	if packet.Header.TotalLength != 20 {
		t.Errorf("Invalid total length: %d", packet.Header.TotalLength)
		return
	}

	if !packet.Header.DestinationIP.Equal(net.ParseIP("192.168.0.1")) {
		t.Error("Invalid source ip")
		return
	}

	if packet.Header.GetProtocol() != ProtocolTCP {
		t.Errorf("Invalid protocol: %d", packet.Header.Protocol)
		return
	}

	if len(packet.Payload) != 0 {
		t.Error("Payload should be empty")
		return
	}
}

func TestParsePacketNoOptionsWithData(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x54, 0xbf, 0x08, 0x40, 0x00, 0x40, 0x01, 0x77, 0xa6, 0x0a, 0xe9, 0xe9, 0x01,
		0x08, 0x08, 0x08, 0x08, 0x08, 0x00, 0xf4, 0x3d, 0xd5, 0x3c, 0x00, 0x01, 0x9a, 0xaf, 0x5a, 0x69,
		0x00, 0x00, 0x00, 0x00, 0x73, 0x98, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x11, 0x12, 0x13,
		0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
		0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33,
		0x34, 0x35, 0x36, 0x37,
	}

	packet, err := ParsePacket(ipPacket)
	if err != nil {
		t.Error("Packet should parsed", err)
		return
	}

	if packet.Header.TotalLength != 84 {
		t.Errorf("Invalid total length: %d", packet.Header.TotalLength)
		return
	}

	if !packet.Header.SourceIP.Equal(net.ParseIP("10.233.233.1")) {
		t.Error("Invalid source ip")
		return
	}

	if packet.Header.GetProtocol() != ProtocolICMP {
		t.Errorf("Invalid protocol: %d", packet.Header.Protocol)
		return
	}

	headerOptions, err := packet.Header.ParseOptions()
	if err != nil {
		t.Error("Header options should parsed", err)
	}

	if len(headerOptions) != 0 {
		t.Error("Header options should be empty")
	}

	if len(packet.Payload) != 64 {
		t.Error("Payload should correct len", len(packet.Payload))
		return
	}
}
