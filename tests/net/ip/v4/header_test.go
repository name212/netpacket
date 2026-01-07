// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"bytes"
	"testing"

	"github.com/name212/netpacket/net/ip/v4"
	"github.com/stretchr/testify/require"
)

func TestParseIPv4HeaderShortData(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00,
	}

	header, err := v4.ParseHeader(ipPacket)
	require.Error(t, err, "should not parse header")
	require.Nil(t, header, "header should be nil")
}

func TestParseIPv4Header(t *testing.T) {
	ipPacket := []byte{
		0x45, 0x00, 0x00, 0x3c,
		0x1c, 0x46, 0x40, 0x00,
		0x40, 0x06, 0xb1, 0xe6,
		0xc0, 0xa8, 0x00, 0x68,
		0xc0, 0xa8, 0x00, 0x01,
	}

	header := parseHeader(t, ipPacket, 60)

	assertSourceAndDestinationAndProto(t, header, "192.168.0.104", v4.ProtocolTCP, "192.168.0.1", "TCP")

	require.Equal(t, 64, header.GetTTL(), "TTL should be 64")
	require.Equal(t, uint16(45542), header.Checksum, "checksum should be 45542")

	flags := header.GetFlags()
	require.False(t, flags.IsEvil, "flags should not be evil")
	require.False(t, flags.MoreFragments, "flags should not be more fragments")
	require.True(t, flags.DontFragment, "flags should be dont fragment")

	assertNoOptions(t, header)
}

func TestParseIPv4Options(t *testing.T) {
	ipPacket := []byte{
		0x49, 0x00, 0x00, 0x28, 0x03, 0x04,
		0x00, 0x00, 0xfe, 0x01, 0xc1, 0xe0,
		0xaf, 0x2d, 0xb0, 0x00, 0x95, 0xab,
		0x7e, 0x0b, 0x01, 0x82, 0x0b, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
	}

	header := parseHeader(t, ipPacket, 40)

	assertSourceAndDestinationAndProto(t, header, "175.45.176.0", v4.ProtocolICMP, "149.171.126.11", "ICMP")

	options, err := header.ParseOptions()
	require.NoError(t, err)

	require.Len(t, options, 2, "options len should be 2")

	firstOption := options[0]

	assertOptionTypeAndLength(t, firstOption, v4.OptionNoOperation, "NOP", 1)
	require.Empty(t, firstOption.GetData(), "data should be empty for NOP")

	secondOption := options[1]

	assertOptionTypeAndLength(t, secondOption, v4.OptionSecurityRIPSO, "SEC", 11)
	require.Len(t, secondOption.GetData(), 9, "data len should be 9 for SEC")

	expectedData := bytes.Repeat([]byte{0x00}, 9)
	require.Equal(t, expectedData, secondOption.GetData(), "data should correct")
}

func assertOptionTypeAndLength(t *testing.T, option v4.Option, tp v4.OptionType, short string, length int) {
	t.Helper()

	require.Equal(t, tp, option.GetType(), "option type should be %d", option.GetType())
	require.Equal(t, length, option.GetLength(), "option length should be %d", length)
	require.Equal(t, short, option.TypeShort(), "option type should be %s", short)
}

func parseHeader(t *testing.T, data []byte, totalLen int) *v4.Header {
	t.Helper()

	header, err := v4.ParseHeader(data)
	require.NoError(t, err)
	require.NotNil(t, header)

	assertHeaderVersionAndTotalLen(t, header, totalLen)

	return header
}
