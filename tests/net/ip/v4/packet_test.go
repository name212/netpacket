// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"testing"

	"github.com/name212/netpacket"
	"github.com/name212/netpacket/net/ip/v4"
	"github.com/name212/netpacket/tests"
	"github.com/name212/netpacket/transport/udp"
	"github.com/stretchr/testify/require"
)

func TestParseIPv4PacketFailIfDataShort(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x3c,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	_, err := v4.ParsePacket(ipPacket)
	require.Error(t, err, "should fail to parse packet")

	payload, err := v4.ExtractPayload(ipPacket)
	// because packet data len valid but header invalid
	require.NoError(t, err, "payload should be extracted")
	require.Nil(t, payload, "payload should be empty")
}

func TestParseIPv4PacketNoData(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x14,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	packet := parsePacket(t, ipPacket, 20, 0)

	assertSourceAndDestinationAndProto(t, packet.GetHeader(), "192.168.0.104", v4.ProtocolTCP, "192.168.0.1", "TCP")

	payload, err := v4.ExtractPayload(ipPacket)
	require.NoError(t, err, "payload should be extracted")
	require.Empty(t, payload, "payload should be empty")

	// AssertStringer Trim \n from expected
	// use \n this for better observability (show in code as string present)
	expectedHeaderString := `
IPv4 Packet:
	Header:
		Source: 192.168.0.104
		Destination: 192.168.0.1
		Protocol: TCP
		TTL: 64
		Header Size: 20
		Packet Size: 20
		Flags:
			Don't Fragment: true
			More Fragments: false
		No options set
		Checksum: 45542
	Is transport: true
	Payload len: 0
`
	tests.AssertStringer(t, packet, expectedHeaderString)
}

func TestParseIPv4PacketNoOptionsWithData(t *testing.T) {
	packet := icmpValidPacket(t)

	expectedPayload := "CAD0PdU8AAGar1ppAAAAAHOYBwAAAAAAEBESExQVFhcYGRobHB0eHyAhIiMkJSYnKCkqKywtLi8wMTIzNDU2Nw=="
	tests.AssertDataAsBase64(t, expectedPayload, packet.GetPayload(), 64)

	header := packet.GetHeader()
	assertSourceAndDestinationAndProto(t, header, "10.233.233.1", v4.ProtocolICMP, "8.8.8.8", "ICMP")
	assertNoOptions(t, header)

	payload, err := v4.ExtractPayload(icmpValidPacketData)
	require.NoError(t, err, "payload should be extracted")
	tests.AssertDataAsBase64(t, expectedPayload, payload, 64)

	// AssertStringer Trim \n from expected
	// use \n this for better observability (show in code as string present)
	expectedHeaderString := `
IPv4 Packet:
	Header:
		Source: 10.233.233.1
		Destination: 8.8.8.8
		Protocol: ICMP
		TTL: 64
		Header Size: 20
		Packet Size: 84
		Flags:
			Don't Fragment: true
			More Fragments: false
		No options set
		Checksum: 30630
	Is transport: false
	Payload len: 64
`
	tests.AssertStringer(t, packet, expectedHeaderString)
}

func TestParseIPv4PacketWithOptionsWithData(t *testing.T) {
	var icmpWithOptions = []byte{
		0x49, 0x00, 0x00, 0x64, 0x03, 0x04,
		0x00, 0x00, 0xfe, 0x01, 0xc1, 0xe0,
		0xaf, 0x2d, 0xb0, 0x00, 0x95, 0xab,
		0x7e, 0x0b, 0x01, 0x82, 0x0b, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x03,

		0x08, 0x00, 0xf4, 0x3d, 0xd5, 0x3c,
		0x00, 0x01, 0x9a, 0xaf, 0x5a, 0x69,
		0x00, 0x00, 0x00, 0x00, 0x73, 0x98,
		0x07, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15,
		0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b,
		0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21,
		0x22, 0x23, 0x24, 0x25, 0x26, 0x27,
		0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d,
		0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33,
		0x34, 0x35, 0x36, 0x37,
	}

	packet := parsePacket(t, icmpWithOptions, 100, 64)

	expectedPayload := "CAD0PdU8AAGar1ppAAAAAHOYBwAAAAAAEBESExQVFhcYGRobHB0eHyAhIiMkJSYnKCkqKywtLi8wMTIzNDU2Nw=="
	tests.AssertDataAsBase64(t, expectedPayload, packet.GetPayload(), 64)

	header := packet.GetHeader()
	assertSourceAndDestinationAndProto(t, header, "175.45.176.0", v4.ProtocolICMP, "149.171.126.11", "ICMP")

	options, err := header.ParseOptions()
	require.NoError(t, err, "options should be parsed")
	require.Len(t, options, 2, "should be 2 options")
	require.Equal(t, v4.OptionSecurityRIPSO, options[1].GetType(), "should second option is SEC RIPSO")

	payload, err := v4.ExtractPayload(icmpValidPacketData)
	require.NoError(t, err, "should parse")
	tests.AssertDataAsBase64(t, expectedPayload, payload, 64)

	// AssertStringer Trim \n from expected
	// use \n this for better observability (show in code as string present)
	expectedHeaderString := `
IPv4 Packet:
	Header:
		Source: 175.45.176.0
		Destination: 149.171.126.11
		Protocol: ICMP
		TTL: 254
		Header Size: 36
		Packet Size: 100
		Flags:
			Don't Fragment: false
			More Fragments: false
		Options:
			Option:
				Type: NOP(1)
				Type description: No Operation
				Full Length: 1
				No data
			Option:
				Type: SEC(130)
				Type description: Security RIPSO
				Full Length: 11
				Hex data:
					0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00
					0x00
		Checksum: 49632
	Is transport: false
	Payload len: 64
`
	tests.AssertStringer(t, packet, expectedHeaderString)
}

func TestGetTransportPacket(t *testing.T) {
	t.Run("Not transport", func(t *testing.T) {
		packet := icmpValidPacket(t)
		_, err := packet.TransportPacket()
		require.Error(t, err, "shouldn't be able to get transport packet")
		require.ErrorIs(t, err, v4.ErrNotTransportPacket, "should not implemented")
		require.Contains(t, err.Error(), "not transport packet")
	})

	t.Run("UDP", func(t *testing.T) {
		ipPacket := []byte{
			0x45, 0x00, 0x00, 0x38, 0x56, 0xaf, 0x40, 0x00, 0x40, 0x11, 0x25, 0xe0, 0xac, 0x11,
			0x00, 0x03, 0x09, 0x09, 0x09, 0x09, 0x99, 0x7a, 0x00, 0x35, 0x00, 0x24, 0xbe, 0x5b,
			0x42, 0x22, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x67,
			0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
		}
		packet := parsePacket(t, ipPacket, 56, 36)
		transport, err := packet.TransportPacket()
		require.NoError(t, err, "transport packet should extracted")

		assertTransport(t, transport, 39290, 53, udp.Kind, 28)

		convertToUDP := func() {
			v4.ToUDP(transport)
		}

		require.NotPanics(t, convertToUDP, "should not be able to convert to UDP")
	})

	t.Run("TCP", func(t *testing.T) {
		ipPacket := []byte{
			0x45, 0x00, 0x00, 0x71, 0xd1, 0x69, 0x40, 0x00, 0x40, 0x06, 0xce, 0xc9, 0x0a, 0xe9, 0xe9, 0x01,
			0xd8, 0x3a, 0xce, 0x2e, 0xa7, 0x9e, 0x00, 0x50, 0x00, 0x4d, 0x6b, 0xcc, 0xb7, 0x16, 0x2a, 0x58,
			0x50, 0x18, 0xfa, 0xf0, 0x50, 0xc5, 0x00, 0x00, 0x47, 0x45, 0x54, 0x20, 0x2f, 0x20, 0x48, 0x54,
			0x54, 0x50, 0x2f, 0x31, 0x2e, 0x31, 0x0d, 0x0a, 0x48, 0x6f, 0x73, 0x74, 0x3a, 0x20, 0x67, 0x6f,
			0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x63, 0x6f, 0x6d, 0x0d, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x2d, 0x41,
			0x67, 0x65, 0x6e, 0x74, 0x3a, 0x20, 0x63, 0x75, 0x72, 0x6c, 0x2f, 0x38, 0x2e, 0x35, 0x2e, 0x30,
			0x0d, 0x0a, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x3a, 0x20, 0x2a, 0x2f, 0x2a, 0x0d, 0x0a, 0x0d,
			0x0a,
		}
		packet := parsePacket(t, ipPacket, 113, 93)
		transport, err := packet.TransportPacket()
		require.NoError(t, err, "transport packet should parsed")

		// todo should implement
		// assertTransport(t, transport, 39290, 53, udp.Kind, 28)

		convertToTCP := func() {
			v4.ToTCP(transport)
		}

		require.NotPanics(t, convertToTCP, "should not be able to convert to UDP")
	})
}

func parsePacket(t *testing.T, data []byte, totalLen int, payloadLen int) *v4.Packet {
	t.Helper()

	packet, err := v4.ParsePacket(data)
	require.NoError(t, err, "failed to parse packet")
	require.NotNil(t, packet, "packet should not be nil")

	header := packet.GetHeader()

	require.NotNil(t, header, "packet header should not be nil")

	assertHeaderVersionAndTotalLen(t, header, totalLen)

	require.Len(t, packet.GetPayload(), payloadLen, "payload len should be %d", payloadLen)

	return packet
}

var icmpValidPacketData = []byte{
	0x45, 0x00, 0x00, 0x54, 0xbf, 0x08, 0x40, 0x00, 0x40, 0x01, 0x77, 0xa6, 0x0a, 0xe9, 0xe9, 0x01,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x00, 0xf4, 0x3d, 0xd5, 0x3c, 0x00, 0x01, 0x9a, 0xaf, 0x5a, 0x69,
	0x00, 0x00, 0x00, 0x00, 0x73, 0x98, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x11, 0x12, 0x13,
	0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, 0x21, 0x22, 0x23,
	0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, 0x31, 0x32, 0x33,
	0x34, 0x35, 0x36, 0x37,
}

func icmpValidPacket(t *testing.T) *v4.Packet {
	return parsePacket(t, icmpValidPacketData, 84, 64)
}

func assertTransport(t *testing.T, transport v4.Transport, srcPort, dstPort int, kind netpacket.Kind, payloadLen int) {
	t.Helper()

	require.Equal(t, srcPort, transport.GetSourcePort(), "source port should be %d", srcPort)
	require.Equal(t, dstPort, transport.GetDestinationPort(), "destination port should be %d", dstPort)
	require.Equal(t, kind, transport.Kind(), "transport kind should be %s", kind)
	require.Len(t, transport.GetPayload(), payloadLen, "payload len should be %d", payloadLen)
}
