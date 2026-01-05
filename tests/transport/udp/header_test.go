// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"testing"

	"github.com/name212/netpacket/transport/udp"
	"github.com/stretchr/testify/require"
)

func TestParseUDPHeaderShortData(t *testing.T) {
	datagram := []byte{
		0xd8, 0x2a,
	}

	header, err := udp.ParseHeader(datagram)

	require.Error(t, err, "should not parse")
	require.Nil(t, header)
}

func TestParseUDPHeader(t *testing.T) {
	datagram := []byte{
		0xd8, 0x2a, 0x00, 0x35, 0x00, 0x25, 0xbe, 0x5c,
	}

	header, err := udp.ParseHeader(datagram)

	require.NoError(t, err, "should parse")
	require.NotNil(t, header)

	assertHeader(t, header, 55338, 53, 37, 48732)
}
