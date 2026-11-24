// Copyright 2026
// license that can be found in the LICENSE file.

package unix

import (
	"fmt"
	"strings"
)

const (
	SocketPrefix = "unix:"
)

func ExtractUnixSocketPath(addr string) (string, error) {
	res := strings.TrimPrefix(addr, SocketPrefix)
	if res == "" {
		return "", fmt.Errorf("unix socket path is empty")
	}

	return res, nil
}

func IsUnix(addr string) bool {
	return strings.HasPrefix(addr, SocketPrefix)
}
