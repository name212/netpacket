// Copyright 2026
// license that can be found in the LICENSE file.

package udp

import (
	"testing"

	"github.com/name212/netpacket/transport/udp"
	"github.com/stretchr/testify/require"
)

func assertHeader(t *testing.T, header *udp.Header, srcPort, dstPort, length, checksum int) {
	t.Helper()

	require.Equal(t, srcPort, header.GetSourcePort(), "source port should be %d", srcPort)
	require.Equal(t, dstPort, header.GetDestinationPort(), "destination port should be %d", dstPort)
	require.Equal(t, uint16(length), header.Length, "length should be %d", length)
	require.Equal(t, uint16(checksum), header.Checksum, "checksum should be %d", checksum)
}
