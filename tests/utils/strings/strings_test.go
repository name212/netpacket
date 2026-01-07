// Copyright 2026
// license that can be found in the LICENSE file.

package strings

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/name212/netpacket/utils/strings"
)

func TestShiftOnTab(t *testing.T) {
	assertOneShift := func(input, expected string) {
		require.Equal(t, expected, strings.ShiftOnTabs(input, 1))
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
		require.Equal(t, expected, strings.ShiftOnTabs(input, 2))
	}

	assertTwoShifts("", "\t\t")
	assertTwoShifts("\n", "\t\t\n")
	assertTwoShifts("some string", "\t\tsome string")
	assertTwoShifts("some\tstring", "\t\tsome\tstring")
	assertTwoShifts("first string\nsecond string", "\t\tfirst string\n\t\tsecond string")
	assertTwoShifts("first string\nsecond string\n", "\t\tfirst string\n\t\tsecond string\n")
	assertTwoShifts("\nfirst string\nsecond string\n\n", "\t\t\n\t\tfirst string\n\t\tsecond string\n\t\t\n")
}

func TestBytesToHexWithWrap(t *testing.T) {
	assertHex := func(t *testing.T, input []byte, lineLen int, expected string) {
		res := strings.BytesToHexWithWrap(input, lineLen)
		require.Equal(t, expected, res, "should format")
	}

	t.Run("zero input", func(t *testing.T) {
		assertHex(t, nil, 2, "")
	})

	t.Run("zero lineLen", func(t *testing.T) {
		b := []byte{0xFF, 0x00, 0x01, 0xA0}
		assertHex(t, b, 0, "0xFF 0x00 0x01 0xA0")
	})

	t.Run("lineLen > len(data)", func(t *testing.T) {
		b := []byte{0xFF, 0x00, 0x01, 0xA0}
		assertHex(t, b, 8, "0xFF 0x00 0x01 0xA0")
	})

	t.Run("lineLen == len(data)", func(t *testing.T) {
		b := []byte{0xFF, 0x00, 0x01, 0xA0}
		assertHex(t, b, 4, "0xFF 0x00 0x01 0xA0")
	})

	t.Run("two lines", func(t *testing.T) {
		b := []byte{
			0xFF, 0x00, 0x01, 0xA0,
			0xAA, 0x11, 0x08, 0xB0,
		}
		assertHex(t, b, 4, `0xFF 0x00 0x01 0xA0
0xAA 0x11 0x08 0xB0`)
	})

	t.Run("three lines", func(t *testing.T) {
		b := []byte{
			0xFF, 0x00, 0x01, 0xA0,
			0xAA, 0x11, 0x08, 0xB0,
			0x76, 0x02, 0x80, 0xB0,
		}
		assertHex(t, b, 4, `0xFF 0x00 0x01 0xA0
0xAA 0x11 0x08 0xB0
0x76 0x02 0x80 0xB0`)
	})

	t.Run("last line is not full", func(t *testing.T) {
		b := []byte{
			0xFF, 0x00, 0x01, 0xA0,
			0xAA, 0x11, 0x08, 0xB0,
			0x76, 0x02, 0xBB,
		}
		assertHex(t, b, 4, `0xFF 0x00 0x01 0xA0
0xAA 0x11 0x08 0xB0
0x76 0x02 0xBB`)
		b = []byte{
			0xFF, 0x00, 0x01, 0xA0,
			0xAA, 0x11, 0x08, 0xB0,
			0x76, 0x02,
		}
		assertHex(t, b, 4, `0xFF 0x00 0x01 0xA0
0xAA 0x11 0x08 0xB0
0x76 0x02`)

		b = []byte{
			0xFF, 0x00, 0x01, 0xA0,
			0xAA, 0x11, 0x08, 0xB0,
			0x76,
		}
		assertHex(t, b, 4, `0xFF 0x00 0x01 0xA0
0xAA 0x11 0x08 0xB0
0x76`)
	})
}
