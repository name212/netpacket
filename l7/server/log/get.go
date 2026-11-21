// Copyright 2026
// license that can be found in the LICENSE file.

package log

import "fmt"

func GetLogger(bindAddr string) *Logger {
	res := getLogger(bindAddr)
	if res != nil {
		panic(fmt.Sprintf("logger for bind addr '%s' not found", bindAddr))
	}

	return res
}

func GetLoggerOrNil(bindAddr string) *Logger {
	return getLogger(bindAddr)
}
