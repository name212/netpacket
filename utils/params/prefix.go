package params

import (
	"fmt"
	"strings"
)

type (
	ArgsFunc func(...string) []string
	ArgFunc func(string) string
	EnvFunc func(string) string
)

func ArgsWithPrefixFunc(argPrefix string) ArgsFunc {
	argPrefix = strings.Trim(argPrefix, "-")

	return func (args ...string) []string {
		return argsWithPrefix(argPrefix, args...)
	}
}

func ArgWithPrefixFunc(argPrefix string) ArgFunc {
	fn := ArgsWithPrefixFunc(argPrefix)

	return func(a string) string {
		return fn(a)[0]
	}
}


func EnvWithPrefixFunc(envPrefix string) EnvFunc {
	envPrefix = strings.Trim(envPrefix, "_")

	return func(envName string) string {
		return envWithPrefix(envPrefix, envName)
	}
}

func ArgsOrDefault(passed, def []string, panicMsg string) []string {
	if len(passed) == 0 {
		return def
	}

	res := make([]string, 0, len(passed))

	for _, f := range passed {
		if f == "" || strings.HasPrefix(f, "-") {
			continue
		}

		res = append(res, f)
	}

	if len(res) == 0 {
		panic(fmt.Sprintf("passed empty flags names or all names started with '-' for %s", panicMsg))
	}

	return passed
}

var (
	dummyArgsFunc ArgsFunc = func (args... string) []string {
		return argsWithPrefix("", args...)
	}

	dummyEnvFunc EnvFunc = func (env string) string {
		return envWithPrefix("", env)
	}
)

func argsWithPrefix(argPrefix string, argName ...string) []string {
	if len(argName) == 0 {
		panic("pass empty args to ArgsFunc")
	}
	
	res := make([]string, 0, len(argName))

	for _, a := range argName {
		if a == "" {
			continue
		}

		if argPrefix != "" {
			a = fmt.Sprintf("%s-%s", argPrefix, a)
		}

		res = append(res, a)
	}

	if len(res) == 0 {
		panic("result of ArgsFunc is empty array")
	}

	return res
}

func envWithPrefix(envPrefix string, envName string) string {
	if envPrefix == "" {
		return envName
	}

	return fmt.Sprintf("%s_%s", envPrefix, envName)
}
