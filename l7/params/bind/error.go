// Copyright 2026
// license that can be found in the LICENSE file.

package bind

import (
	"fmt"
	"strings"
)

func errMsgForArgsEnv(suf string, args []string, env string) string {
	return fmt.Sprintf(
		"%s as args [%s] or env %s",
		suf,
		strings.Join(args, ","),
		env,
	)
}

func checkErrValidate(suf string, err error) error {
	if err != nil {
		return fmt.Errorf("%s is incorrect: %w", suf, err)
	}

	return nil
}

func appendErrForBind(err error, f string, args ...any) error {
	curErr := fmt.Errorf(f, args...)
	if err != nil {
		return fmt.Errorf("%v\n%v", err, curErr)
	}

	return curErr
}

func appendErrForBindAsStr(err error, s string) error {
	return appendErrForBind(err, "%s", s)
}
