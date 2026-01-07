// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"testing"

	"github.com/name212/netpacket/tests"
	"github.com/name212/netpacket/transport/udp"
	"github.com/stretchr/testify/require"
)

func TestParseUDPDatagramShortData(t *testing.T) {
	data := []byte{
		0xd8, 0x2a,
	}

	datagram, err := udp.ParseDatagram(data)
	require.Error(t, err, "should not parse")
	require.Nil(t, datagram)

	payload, err := udp.ExtractPayload(data)
	require.Error(t, err, "should not extract payload")
	require.Nil(t, payload, "should extract empty payload")
}

func TestParseUDPDatagramWithoutPayload(t *testing.T) {
	data := []byte{
		0xd8, 0x2a, 0x00, 0x35, 0x00, 0x25, 0xbe, 0x5c,
	}

	datagram := parseDatagram(t, data)

	assertHeader(t, datagram.GetHeader(), 55338, 53, 37, 48732)
	require.Empty(t, datagram.GetPayload(), "should parse empty payload")

	payload, err := udp.ExtractPayload(data)
	require.NoError(t, err, "should extract payload")
	require.Empty(t, payload, "should extract empty payload")
}

func TestParseUDPDatagramWithPayload(t *testing.T) {
	data := []byte{
		0x99, 0x7a, 0x00, 0x35, 0x00, 0x24, 0xbe, 0x5b,
		0x42, 0x22, 0x01, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x67, 0x6f, 0x6f,
		0x67, 0x6c, 0x65, 0x03, 0x63, 0x6f, 0x6d, 0x00, 0x00, 0x01, 0x00, 0x01,
	}

	const payloadLength = 28

	datagram := parseDatagram(t, data)

	assertHeader(t, datagram.GetHeader(), 39290, 53, 36, 48731)

	expectedPayload := "QiIBAAABAAAAAAAABmdvb2dsZQNjb20AAAEAAQ=="

	tests.AssertDataAsBase64(t, expectedPayload, datagram.GetPayload(), payloadLength)

	payload, err := udp.ExtractPayload(data)
	require.NoError(t, err, "should extract payload")
	tests.AssertDataAsBase64(t, expectedPayload, payload, payloadLength)
}

func parseDatagram(t *testing.T, data []byte) *udp.Datagram {
	datagram, err := udp.ParseDatagram(data)
	require.NoError(t, err, "should parse")
	require.NotNil(t, datagram, "should parse")

	return datagram
}
