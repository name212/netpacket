// Copyright 2026
// license that can be found in the LICENSE file.

package bind

import (
	"fmt"

	"github.com/name212/netpacket/utils/params"
)

type AddrParams struct {
	HelpSuffix string

	FullBindAddrArgName []string
	FullBindAddrEnvName string
	FullBindAddrDefault string

	BindHostArgName []string
	BindHostEnvName string
	BindHostDefault string

	BindPortArgName []string
	BindPortEnvName string
	BindPortDefault string

	UnixSocketModeArgName []string
	UnixSocketModeEnvName string

	argsFunc params.ArgsFunc
	envFunc  params.EnvFunc
}

func (p *AddrParams) help(h string) string {
	if p.HelpSuffix == "" {
		return h
	}

	return fmt.Sprintf("%s for %s", h, p.HelpSuffix)
}

func (p *AddrParams) errMsgPort() string {
	return errMsgForArgsEnv(portMsgSuf, p.BindPortArgName, p.BindPortEnvName)
}

func (p *AddrParams) errMsgHost() string {
	return errMsgForArgsEnv(hostMsgSuf, p.BindHostArgName, p.BindHostEnvName)
}

func (p *AddrParams) errMsgFull() string {
	return errMsgForArgsEnv(fullAddrMsgSuf, p.FullBindAddrArgName, p.FullBindAddrEnvName)
}

func bindParamOrDefault(passed, def string) string {
	if passed == "" {
		return def
	}

	return passed
}

func getBindParamWithDefaults(p *AddrParams) (*AddrParams, error) {
	defaults := DefaultBindAddrParams()
	res := &AddrParams{}

	res.HelpSuffix = bindParamOrDefault(p.HelpSuffix, "")

	res.FullBindAddrArgName = params.ArgsOrDefault(p.FullBindAddrArgName, defaults.FullBindAddrArgName, fullAddrMsgSuf)
	res.FullBindAddrEnvName = bindParamOrDefault(p.FullBindAddrEnvName, defaults.FullBindAddrEnvName)
	res.FullBindAddrDefault = bindParamOrDefault(p.FullBindAddrDefault, "")

	res.BindHostArgName = params.ArgsOrDefault(p.BindHostArgName, defaults.BindHostArgName, hostMsgSuf)
	res.BindHostEnvName = bindParamOrDefault(p.BindHostEnvName, defaults.BindHostEnvName)
	res.BindHostDefault = bindParamOrDefault(p.BindHostDefault, "")

	res.BindPortArgName = params.ArgsOrDefault(p.BindPortArgName, defaults.BindPortArgName, portMsgSuf)
	res.BindPortEnvName = bindParamOrDefault(p.BindPortEnvName, defaults.BindPortEnvName)
	res.BindPortDefault = bindParamOrDefault(p.BindPortDefault, "")

	res.UnixSocketModeArgName = params.ArgsOrDefault(p.UnixSocketModeArgName, defaults.UnixSocketModeArgName, unixSocketModeMsgSuf)
	res.UnixSocketModeEnvName = bindParamOrDefault(p.UnixSocketModeEnvName, defaults.UnixSocketModeEnvName)

	var errDefaults error

	validate := func(v string, vFunc func(string) error) {
		if err := vFunc(v); err != nil {
			errDefaults = appendErrForBind(errDefaults, "default val for %s", err.Error())
		}
	}

	if res.FullBindAddrDefault != "" {
		validate(res.FullBindAddrDefault, ValidateBindFullAddr)
	}

	if res.BindHostDefault != "" {
		validate(res.BindHostDefault, ValidateBindHost)
	}

	if res.BindPortDefault != "" {
		validate(res.BindPortDefault, ValidateBindPort)
	}

	if errDefaults != nil {
		return nil, fmt.Errorf("%s:\n%s", p.help("default values incorrect"), errDefaults.Error())
	}

	return res, nil
}
