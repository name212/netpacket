// Copyright 2026
// license that can be found in the LICENSE file.

package bind

import (
	"flag"
	"fmt"
	"net"

	"github.com/name212/netpacket/utils/params"
	"github.com/name212/netpacket/l7/params/unix"
)

const (
	DefaultFullBindAddrArg     = "bind-addr"
	DefaultFullBindAddrEnvName = "BIND_ADDR"

	DefaultHostArg     = "bind-host"
	DefaultHostEnvName = "BIND_HOST"

	DefaultPortArgName = "bind-port"
	DefaultPortEnvName = "BIND_PORT"

	DefaultUnixSocketModeArgName = "unix-socket-mode"
	DefaultUnixSocketModeEnvName = "UNIX_SOCKET_MODE"

	defaultUnixSocketModeInt = 0o700
	DefaultUnixSocketMode    = "700"
)

type Addr struct {
	BindAddr          string
	UnixSocketModeStr string
}

func NewAddr(bindAddr string) *Addr {
	return &Addr{
		BindAddr: bindAddr,
	}
}

func DefaultBindAddrParams() *AddrParams {
	return BindAddrParamsWithPrefixes("", "")
}

func BindAddrParamsWithPrefixes(argPrefix, envPrefix string) *AddrParams {
	var argsFunc params.ArgsFunc
	if argPrefix != "" {
		argsFunc = params.ArgsWithPrefixFunc(argPrefix)
	}

	var envFunc params.EnvFunc
	if envPrefix != "" {
		envFunc = params.EnvWithPrefixFunc(envPrefix)
	}

	return &AddrParams{
		FullBindAddrArgName: []string{DefaultFullBindAddrArg},
		FullBindAddrEnvName: DefaultFullBindAddrEnvName,

		BindHostArgName: []string{DefaultHostArg},
		BindHostEnvName: DefaultHostEnvName,

		BindPortArgName: []string{DefaultPortArgName},
		BindPortEnvName: DefaultPortEnvName,

		UnixSocketModeArgName: []string{DefaultUnixSocketModeArgName},
		UnixSocketModeEnvName: DefaultUnixSocketModeEnvName,

		argsFunc: argsFunc,
		envFunc:  envFunc,
	}
}

func InitBindFlags(p *AddrParams, set ...*flag.FlagSet) (*AddrFlags, error) {
	if p == nil {
		return nil, fmt.Errorf("pass nil AddrParams to InitBindFlags")
	}

	p, err := getBindParamWithDefaults(p)
	if err != nil {
		return nil, err
	}

	return newAddrFlags(p, set...)
}

func GetBindAddress(f *AddrFlags) (*Addr, error) {
	if err := f.isInitialized(); err != nil {
		return nil, err
	}

	bindFullAddr, err := f.extractFull()
	if err != nil {
		return nil, err
	}

	bindHost, err := f.extractHost()
	if err != nil {
		return nil, err
	}

	bindPort, err := f.extractPort()
	if err != nil {
		return nil, err
	}

	if bindFullAddr != "" {
		if err := ValidateBindFullAddr(bindFullAddr); err != nil {
			return nil, fmt.Errorf("%s", f.params.help(err.Error()))
		}

		res := &Addr{
			BindAddr: bindFullAddr,
		}

		if res.IsUnix() {
			var errUnix error

			if _, err := res.ExtractUnixSocketPath(); err != nil {
				errUnix = appendErrForBindAsStr(errUnix, err.Error())
			}

			unixModeStr, err := f.extractUnixSocketMode()
			if err == nil {
				if err := ValidateUnixSocketMode(unixModeStr); err != nil {
					errUnix = appendErrForBindAsStr(errUnix, err.Error())
				}
			} else {
				errUnix = appendErrForBindAsStr(errUnix, err.Error())
			}

			if errUnix != nil {
				return nil, fmt.Errorf(
					"%s %s:\n%s",
					f.params.errMsgFull(),
					f.params.help("unix socket params"),
					errUnix.Error(),
				)
			}

			res.UnixSocketModeStr = unixModeStr
		}

		return res, nil
	}

	var errPortHostEmpty error

	if bindHost == "" {
		errPortHostEmpty = appendErrForBindAsStr(errPortHostEmpty, f.params.errMsgHost())
	}

	if bindPort == "" {
		errPortHostEmpty = appendErrForBindAsStr(errPortHostEmpty, f.params.errMsgPort())
	}

	if errPortHostEmpty != nil {
		return nil, fmt.Errorf(
			"%s %s:\n%s",
			f.params.errMsgFull(),
			f.params.help("not passed and next args not passed"),
			errPortHostEmpty.Error(),
		)
	}

	var errValidate error

	if err := ValidateBindHost(bindHost); err != nil {
		errValidate = appendErrForBindAsStr(errValidate, err.Error())
	}

	if err := ValidateBindPort(bindPort); err != nil {
		errValidate = appendErrForBindAsStr(errValidate, err.Error())
	}

	if errValidate != nil {
		return nil, fmt.Errorf(
			"%s:\n%s",
			f.params.help("cannot validate"),
			errValidate.Error(),
		)
	}

	return &Addr{
		BindAddr: net.JoinHostPort(bindHost, bindPort),
	}, nil
}

func (a *Addr) Addr() string {
	return a.BindAddr
}

func (a *Addr) UnixSocketMode() (int, error) {
	if a.UnixSocketModeStr == "" {
		return defaultUnixSocketModeInt, nil
	}

	return validateUnixSocketMode(a.UnixSocketModeStr)
}

func (a *Addr) ExtractUnixSocketPath() (string, error) {
	return unix.ExtractUnixSocketPath(a.BindAddr)
}

func (a *Addr) IsUnix() bool {
	return unix.IsUnix(a.BindAddr)
}
