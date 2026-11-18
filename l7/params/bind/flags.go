package bind

import (
	"flag"
	"fmt"

	"github.com/name212/netpacket/utils/params"
)

type AddrFlags struct {
	fullBindAddrParam *params.Param
	bindHostParam     *params.Param
	bindPortParam     *params.Param

	params *AddrParams
}

func newAddrFlags(p *AddrParams, set ...*flag.FlagSet) (*AddrFlags, error) {
	res := &AddrFlags{
		fullBindAddrParam: params.NewParam(
			fullAddrMsgSuf,
			p.FullBindAddrArgName,
			p.FullBindAddrEnvName,
			p.help("full bind address 'host:port'"),
		).WithDefault(p.FullBindAddrDefault),

		bindHostParam: params.NewParam(
			hostMsgSuf,
			p.BindHostArgName,
			p.BindHostEnvName,
			p.help("bind host address"),
		).WithDefault(p.BindHostDefault),

		bindPortParam: params.NewParam(
			portMsgSuf,
			p.BindPortArgName,
			p.BindPortEnvName,
			p.help("bind port"),
		).WithDefault(p.BindPortDefault),
		params: p,
	}

	params.SetFuncsForParams(p.argsFunc, p.envFunc, res.fullBindAddrParam, res.bindHostParam, res.bindPortParam)

	if err := res.fullBindAddrParam.Init(set...); err != nil {
		return nil, err
	}

	if err := res.bindHostParam.Init(set...); err != nil {
		return nil, err
	}

	if err := res.bindPortParam.Init(); err != nil {
		return nil, err
	}

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
	return f.fullBindAddrParam.ExtractString()
}

func (f *AddrFlags) extractHost() (string, error) {
	return f.bindHostParam.ExtractString()
}

func (f *AddrFlags) extractPort() (string, error) {
	return f.bindPortParam.ExtractString()
}
