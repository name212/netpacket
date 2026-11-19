package bind

import (
	"flag"
	"fmt"
	"net"

	"github.com/name212/netpacket/utils/params"
)

const (
	DefaultFullBindAddrArg     = "bind-addr"
	DefaultFullBindAddrEnvName = "BIND_ADDR"

	DefaultHostArg     = "bind-host"
	DefaultHostEnvName = "BIND_HOST"

	DefaultPortArgName = "bind-port"
	DefaultPortEnvName = "BIND_PORT"
)

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

func GetBindAddress(f *AddrFlags) (string, error) {
	if err := f.isInitialized(); err != nil {
		return "", err
	}

	bindFullAddr, err := f.extractFull()
	if err != nil {
		return "", err
	}

	bindHost, err := f.extractHost()
	if err != nil {
		return "", err
	}

	bindPort, err := f.extractPort()
	if err != nil {
		return "", err
	}

	if bindFullAddr != "" {
		if err := ValidateBindFullAddr(bindFullAddr); err != nil {
			return "", fmt.Errorf("%s", f.params.help(err.Error()))
		}

		return bindFullAddr, nil
	}

	var errPortHostEmpty error

	if bindHost == "" {
		errPortHostEmpty = appendErrForBindAsStr(errPortHostEmpty, f.params.errMsgHost())
	}

	if bindPort == "" {
		errPortHostEmpty = appendErrForBindAsStr(errPortHostEmpty, f.params.errMsgPort())
	}

	if errPortHostEmpty != nil {
		return "", fmt.Errorf(
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
		return "", fmt.Errorf(
			"%s:\n%s",
			f.params.help("cannot validate"),
			errValidate.Error(),
		)
	}

	return net.JoinHostPort(bindHost, bindPort), nil
}
