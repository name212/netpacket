// Copyright 2026
// license that can be found in the LICENSE file.

package utils

import (
	"testing"

	"github.com/name212/netpacket/utils"
	"github.com/stretchr/testify/require"
)

func TestShiftOnTab(t *testing.T) {
	assertOneShift := func(input, expected string) {
		require.Equal(t, expected, utils.ShiftOnTabs(input, 1))
	}

	assertOneShift("", "\t")
	assertOneShift("\n", "\t\n")
	assertOneShift("some string", "\tsome string")
	assertOneShift("some\tstring", "\tsome\tstring")
	assertOneShift("first string\nsecond string", "\tfirst string\n\tsecond string")
	assertOneShift("first string\nsecond string\n", "\tfirst string\n\tsecond string\n")
	assertOneShift("\nfirst string\nsecond string\n\n", "\t\n\tfirst string\n\tsecond string\n\t\n")

}

func TestMultipleShiftsOnTab(t *testing.T) {
	assertTwoShifts := func(input, expected string) {
		require.Equal(t, expected, utils.ShiftOnTabs(input, 2))
	}

	assertTwoShifts("", "\t\t")
	assertTwoShifts("\n", "\t\t\n")
	assertTwoShifts("some string", "\t\tsome string")
	assertTwoShifts("some\tstring", "\t\tsome\tstring")
	assertTwoShifts("first string\nsecond string", "\t\tfirst string\n\t\tsecond string")
	assertTwoShifts("first string\nsecond string\n", "\t\tfirst string\n\t\tsecond string\n")
	assertTwoShifts("\nfirst string\nsecond string\n\n", "\t\t\n\t\tfirst string\n\t\tsecond string\n\t\t\n")
}
