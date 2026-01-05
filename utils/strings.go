// Copyright 2026
// license that can be found in the LICENSE file.

package utils

import "fmt"

func FmtLn(f string, args ...any) string {
	f = f + "\n"
	return fmt.Sprintf(f, args...)
}

func FmtLnWithTabPrefix(f string, args ...any) string {
	f = "\t" + f + "\n"
	return fmt.Sprintf(f, args...)
}
