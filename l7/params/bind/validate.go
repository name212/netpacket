package bind

import (
	"fmt"
	"net"
	"net/netip"
	"strconv"
)

func ValidateBindFullAddr(addr string) error {
	_, _, err := net.SplitHostPort(addr)
	return checkErrValidate(fullAddrMsgSuf, err)
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
