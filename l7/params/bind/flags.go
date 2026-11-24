// Copyright 2026
// license that can be found in the LICENSE file.

package bind

import (
	"flag"
	"fmt"

	"github.com/name212/netpacket/l7/params/unix"
	"github.com/name212/netpacket/utils/params"
)

const (
	fullAddrMsgSuf       = "full bind name"
	hostMsgSuf           = "bind host"
	portMsgSuf           = "bind port"
	unixSocketModeMsgSuf = "unix socket mode"
)

type AddrFlags struct {
	fullBindAddrParam   *params.Param
	bindHostParam       *params.Param
	bindPortParam       *params.Param
	unixSocketModeParam *params.Param

	params *AddrParams
}

func newAddrFlags(p *AddrParams, set ...*flag.FlagSet) (*AddrFlags, error) {
	paramsForInit := make([]*params.Param, 0, 4)

	res := &AddrFlags{}

	addParam := func(target **params.Param, p *params.Param) {
		*target = p
		paramsForInit = append(paramsForInit, p)
	}

	addParam(
		&res.fullBindAddrParam,
		params.NewParam(
			fullAddrMsgSuf,
			p.FullBindAddrArgName,
			p.FullBindAddrEnvName,
			p.help(
				fmt.Sprintf(
					"full bind address 'host:port' or %sPATH for bind to unix domain socket", unix.SocketPrefix,
				),
			),
		).WithDefault(p.FullBindAddrDefault),
	)

	addParam(
		&res.bindHostParam,
		params.NewParam(
			hostMsgSuf,
			p.BindHostArgName,
			p.BindHostEnvName,
			p.help("bind host address"),
		).WithDefault(p.BindHostDefault),
	)

	addParam(
		&res.bindPortParam,
		params.NewParam(
			portMsgSuf,
			p.BindPortArgName,
			p.BindPortEnvName,
			p.help("bind port"),
		).WithDefault(p.BindPortDefault),
	)

	addParam(
		&res.unixSocketModeParam,
		params.NewParam(
			unixSocketModeMsgSuf,
			p.UnixSocketModeArgName,
			p.UnixSocketModeEnvName,
			p.help(
				"if bind full address passed as unix socket with this param you can pass unix file mode (in base 8 number)",
			),
		).WithDefault(DefaultUnixSocketMode),
	)

	params.SetFuncsForParams(p.argsFunc, p.envFunc, paramsForInit...)

	for _, pi := range paramsForInit {
		if err := pi.Init(set...); err != nil {
			return nil, fmt.Errorf(
				"%s: %w",
				p.help("cannot init flags"),
				err,
			)
		}
	}

	res.params = p

	return res, nil
}

func (f *AddrFlags) isInitialized() error {
	if f == nil {
		return fmt.Errorf("AddrFlags is nil")
	}

	if f.params == nil {
		return fmt.Errorf("AddrFlags is not initialized. Please prepare AddrFlags before get")
	}

	return nil
}

func (f *AddrFlags) extractFull() (string, error) {
	return f.extractStringFromParam(f.fullBindAddrParam, fullAddrMsgSuf)
}

func (f *AddrFlags) extractHost() (string, error) {
	return f.extractStringFromParam(f.bindHostParam, hostMsgSuf)
}

func (f *AddrFlags) extractPort() (string, error) {
	return f.extractStringFromParam(f.bindPortParam, portMsgSuf)
}

func (f *AddrFlags) extractUnixSocketMode() (string, error) {
	return f.extractStringFromParam(f.unixSocketModeParam, unixSocketModeMsgSuf)
}

func (f *AddrFlags) extractStringFromParam(p *params.Param, suf string) (string, error) {
	res, err := p.ExtractString()
	if err != nil {
		return "", fmt.Errorf(
			"cannot extract string value %s: %w",
			f.params.help(suf),
			err,
		)
	}

	return res, nil
}
