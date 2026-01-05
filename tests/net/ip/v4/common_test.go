// Copyright 2026
// license that can be found in the LICENSE file.

package v4

import (
	"testing"

	"github.com/name212/netpacket/net/ip/v4"
	"github.com/stretchr/testify/require"
)

func assertSourceAndDestinationAndProto(t *testing.T, header *v4.Header, source string, protocol v4.Protocol, destination string, protocolStr string) {
	t.Helper()

	require.Equal(t, source, header.SourceIP.String(), "source ip should be %s", source)
	require.Equal(t, destination, header.DestinationIP.String(), "destination ip should be %s", destination)
	require.Equal(t, protocol, header.GetProtocol(), "protocol should be %d", protocol)
	require.Equal(t, protocolStr, header.ProtocolString(), "protocol string should be %s", protocolStr)
}

func assertNoOptions(t *testing.T, header *v4.Header) {
	t.Helper()

	headerOptions, err := header.ParseOptions()
	require.NoError(t, err)
	require.Len(t, headerOptions, 0, "header options should be 0 len")
}

func assertHeaderVersionAndTotalLen(t *testing.T, header *v4.Header, length uint16) {
	t.Helper()

	require.Equal(t, uint8(4), header.Version, "version should be 4")
	require.True(t, header.IsValidVersion(), "version should valid")
	require.Equal(t, length, header.TotalLength, "header total length should be %d", length)
}
