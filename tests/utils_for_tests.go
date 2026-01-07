// Copyright 2026
// license that can be found in the LICENSE file.

package tests

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertDataAsBase64(t *testing.T, expected string, data []byte, length int) {
	require.Len(t, data, length, "data len should be %d", length)
	require.Equal(t, expected, base64.StdEncoding.EncodeToString(data), "data should be equal")
}
