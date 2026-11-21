// Copyright 2026
// license that can be found in the LICENSE file.

package params

import (
	"flag"
	"fmt"
	"os"
)

type Param struct {
	set *flagSet

	name string
	help string

	required bool
	def      string

	envName   string
	argsNames []string

	argsFunc ArgsFunc
	envFunc  EnvFunc

	dest []*string
}

func NewParam(name string, argsNames []string, envName string, help string) *Param {
	return &Param{
		name: name,

		help: help,

		envName:   envName,
		argsNames: argsNames,

		argsFunc: dummyArgsFunc,
		envFunc:  dummyEnvFunc,

		dest: make([]*string, 0, len(argsNames)),
	}
}

func SetFuncsForParams(argsFunc ArgsFunc, envFunc EnvFunc, params ...*Param) {
	for _, p := range params {
		p.WithArgsFunc(argsFunc).WithEnvFunc(envFunc)
	}
}

func ExtractParam[T any](p *Param, converter func(string) (T, error)) (T, error) {
	var zero T

	strVal, err := p.ExtractString()
	if err != nil {
		return zero, err
	}

	return converter(strVal)
}

func (p *Param) Init(set ...*flag.FlagSet) error {
	s, err := newFlagSet(set)
	if err != nil {
		return err
	}

	p.set = s

	names := p.argsFunc(p.argsNames...)

	defMsg := ""
	if p.def != "" {
		defMsg = ". By default " + p.def
	}

	help := fmt.Sprintf("%s (can be set with env %s)%s", p.help, p.getEnvName(), defMsg)

	return p.set.initFlags(names, &p.dest, help)
}

func (p *Param) WithArgsFunc(f ArgsFunc) *Param {
	if f != nil {
		p.argsFunc = f
	}

	return p
}

func (p *Param) WithEnvFunc(f EnvFunc) *Param {
	if f != nil {
		p.envFunc = f
	}

	return p
}

func (p *Param) WithRequired(f bool) *Param {
	p.required = f
	return p
}

func (p *Param) WithDefault(d string) *Param {
	p.def = d
	return p
}

func (p *Param) ExtractString() (string, error) {
	if err := p.isInit(); err != nil {
		return "", err
	}

	if err := p.set.NotParsed(); err != nil {
		return "", err
	}

	envName := p.getEnvName()

	if envVal := os.Getenv(envName); envVal != "" {
		return envVal, nil
	}

	res := ""

	for _, vPtr := range p.dest {
		v := *vPtr

		if v == "" {
			continue
		}

		if res != "" {
			return "", fmt.Errorf("found multiple flags passed for Param '%s'", p.name)
		}

		res = v
	}

	if res == "" {
		if p.required {
			return "", fmt.Errorf("value for Param '%s' not found", p.name)
		}

		res = p.def

		if res == "" && p.required {
			return "", fmt.Errorf("value for Param '%s' not found and default val not set", p.name)
		}
	}

	return res, nil
}

func (p *Param) getEnvName() string {
	return p.envFunc(p.envName)
}

func (p *Param) isInit() error {
	if p == nil {
		return fmt.Errorf("Param is nil")
	}

	if p.set == nil {
		return fmt.Errorf("Param '%s' is not initialized. Please call Init before extract", p.name)
	}

	if len(p.dest) == 0 {
		return fmt.Errorf("Param '%s' is not initialized, empty destination array", p.name)
	}

	for i := range p.dest {
		if p.dest[i] == nil {
			return fmt.Errorf("Param '%s' is not initialized, nil string at %d", p.name, i)
		}
	}

	return nil
}
