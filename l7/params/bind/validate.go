// Copyright 2026
// license that can be found in the LICENSE file.

package bind

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
)

func ValidateBindFullAddr(addr string) error {
	a := NewAddr(addr)

	if a.IsUnix() {
		_, err := a.ExtractUnixSocketPath()
		return checkErrValidate(fullAddrMsgSuf, err)
	}

	host, port, err := net.SplitHostPort(a.BindAddr)
	if err := checkErrValidate(fullAddrMsgSuf, err); err != nil {
		return err
	}

	if err := ValidateBindHost(host); err != nil {
		return checkErrValidate(fullAddrMsgSuf, err)
	}

	if err := ValidateBindPort(port); err != nil {
		return checkErrValidate(fullAddrMsgSuf, err)
	}

	return nil
}

func ValidateBindHost(bindHost string) error {
	_, err := netip.ParseAddr(bindHost)
	return checkErrValidate(hostMsgSuf, err)
}

func ValidateBindPort(bindPort string) error {
	portInt, err := strconv.ParseUint(bindPort, 10, 16)
	if err != nil {
		return checkErrValidate(portMsgSuf, fmt.Errorf("not number '%w'", err))
	}

	if portInt < 1 || portInt > 65535 {
		return checkErrValidate(portMsgSuf, fmt.Errorf("should > 0 and <= 65535"))
	}

	return nil
}

const maxUnixMode = 0o777

func ValidateUnixSocketMode(mode string) error {
	_, err := validateUnixSocketMode(mode)
	return err
}

func validateUnixSocketMode(mode string) (int, error) {
	errVal := -1
	modeInt, err := strconv.ParseUint(mode, 8, 16)
	if err != nil {
		return errVal, checkErrValidate(unixSocketModeMsgSuf, fmt.Errorf("not number '%w'", err))
	}

	if modeInt < 1 || modeInt > maxUnixMode {
		return errVal, checkErrValidate(
			unixSocketModeMsgSuf,
			fmt.Errorf(
				"should > 0 and <= %d (%o in base 8)",
				maxUnixMode,
				maxUnixMode,
			),
		)
	}

	return int(modeInt), nil
}
