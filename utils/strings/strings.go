// Copyright 2026
// license that can be found in the LICENSE file.

package strings

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

// FmtWithTabPrefix
// Add one tab before format and Sprintf format with arguments
func FmtWithTabPrefix(f string, args ...any) string {
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

// BytesToHexWithWrap
// returns hex representation of string
// split by new lines each line has stringLen
func BytesToHexWithWrap(data []byte, lineLen int) string {
	dataLen := len(data)
	if dataLen == 0 {
		return ""
	}

	formatLine := func(line []byte) string {
		lineBuilder := strings.Builder{}

		// 4 char for byte representation like 0xAA
		// and spaces between bytes representation
		lineBuilder.Grow(lineLen*4 + len(data))

		for _, ch := range line {
			lineBuilder.WriteString(fmt.Sprintf("0x%02X ", ch))
		}

		return strings.TrimSuffix(lineBuilder.String(), " ")
	}

	if lineLen < 1 || dataLen <= lineLen {
		return formatLine(data)
	}

	b := strings.Builder{}
	// 4 char for byte representation like 0xAA and space separator plus new line
	// it is approximately
	b.Grow(dataLen * 6)

	d := data
	for len(d) >= lineLen {
		b.WriteString(formatLine(d[0:lineLen]))
		b.WriteString("\n")
		d = d[lineLen:]
	}

	if len(d) > 0 {
		b.WriteString(formatLine(d))
	}

	return strings.TrimSuffix(b.String(), "\n")
}

func tabsStr(count int) string {
	return strings.Repeat("\t", count)
}
