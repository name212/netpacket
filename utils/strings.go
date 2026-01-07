// Copyright 2026
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"strings"
)

// FmtLn
// Add new line to format string and Sprintf format with arguments
func FmtLn(f string, args ...any) string {
	f += "\n"
	return fmt.Sprintf(f, args...)
}

// FmtLnWithTabPrefix
// Add one tab before format and add new line after format
// and Sprintf format with arguments
func FmtLnWithTabPrefix(f string, args ...any) string {
	f = tabsStr(1) + f + "\n"
	return fmt.Sprintf(f, args...)
}

// FmtWithTabsPrefix
// Add one tab before format and Sprintf format with arguments
func FmtWithTabsPrefix(f string, args ...any) string {
	f = tabsStr(1) + f
	return fmt.Sprintf(f, args...)
}

// ShiftOnTabs
// Add more tabs in all lines in string
// If last string contains new line it will save it new line
// without add tabs
func ShiftOnTabs(s string, tabsCount int) string {
	tabs := tabsStr(tabsCount)
	lnWithTabs := fmt.Sprintf("\n%s", tabs)
	res := strings.ReplaceAll(s, "\n", lnWithTabs)
	if strings.HasSuffix(res, lnWithTabs) {
		res = fmt.Sprintf("%s\n", strings.TrimSuffix(res, lnWithTabs))
	}

	return tabs + res
}

func tabsStr(count int) string {
	return strings.Repeat("\t", count)
}
