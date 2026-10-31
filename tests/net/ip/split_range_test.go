// Copyright 2026
// license that can be found in the LICENSE file.

package ip

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/name212/netpacket/net/ip"
)

func TestSplitRangeErrors(t *testing.T) {
	ipV4Start := mustParseAddress(t, "192.168.1.22")
	ipV4End := mustParseAddress(t, "192.168.1.30")
	ipV6Start := mustParseAddress(t, "fe80:0000:0000:0000:0202:b3ff:fe1e:8329")
	ipV6End := mustParseAddress(t, "fe80:0000:0000:0000:0202:b3ff:fe1e:ff00")

	t.Run("first v4 second v6", func(t *testing.T) {
		_, err := ip.SplitRangeToSubnets(ipV4Start, ipV6Start)
		require.Error(t, err, "should returns error")
	})

	t.Run("first v6 second v4", func(t *testing.T) {
		_, err := ip.SplitRangeToSubnets(ipV6Start, ipV4Start)
		require.Error(t, err, "should returns error")
	})

	t.Run("v4 incorrect range", func(t *testing.T) {
		_, err := ip.SplitRangeToSubnets(ipV4End, ipV4Start)
		require.Error(t, err, "should returns error")
	})

	t.Run("v6 incorrect range", func(t *testing.T) {
		_, err := ip.SplitRangeToSubnets(ipV6End, ipV6Start)
		require.Error(t, err, "should returns error")
	})
}

type splitRangeToSubnetsTest struct {
	name     string
	start    string
	end      string
	expected []string
}

func TestSplitRangeV4(t *testing.T) {
	cases := []splitRangeToSubnetsTest{
		{
			name:  "one address",
			start: "192.168.0.1",
			end:   "192.168.0.1",
			expected: []string{
				"192.168.0.1/32",
			},
		},

		{
			name:  "two addresses",
			start: "192.168.0.1",
			end:   "192.168.0.2",
			expected: []string{
				"192.168.0.1/32",
				"192.168.0.2/32",
			},
		},

		{
			name:  "three addresses",
			start: "192.168.0.1",
			end:   "192.168.0.3",
			expected: []string{
				"192.168.0.1/32",
				"192.168.0.2/31",
			},
		},

		{
			name:  "four addresses",
			start: "192.168.0.1",
			end:   "192.168.0.4",
			expected: []string{
				"192.168.0.1/32",
				"192.168.0.2/31",
				"192.168.0.4/32",
			},
		},

		{
			name:  "over one subnet from zero",
			start: "192.168.1.0",
			end:   "192.168.1.20",
			expected: []string{
				"192.168.1.0/28",
				"192.168.1.16/30",
				"192.168.1.20/32",
			},
		},

		{
			name:  "over one subnet from one",
			start: "192.168.1.1",
			end:   "192.168.1.20",
			expected: []string{
				"192.168.1.1/32",
				"192.168.1.2/31",
				"192.168.1.4/30",
				"192.168.1.8/29",
				"192.168.1.16/30",
				"192.168.1.20/32",
			},
		},

		{
			name:  "over one subnet in middle not align",
			start: "192.168.1.10",
			end:   "192.168.1.23",
			expected: []string{
				"192.168.1.10/31",
				"192.168.1.12/30",
				"192.168.1.16/29",
			},
		},

		{
			name:  "over one subnet in middle with align",
			start: "192.168.1.16",
			end:   "192.168.1.31",
			expected: []string{
				"192.168.1.16/28",
			},
		},

		{
			name:  "over one subnet 24 all subnet",
			start: "192.168.1.0",
			end:   "192.168.1.255",
			expected: []string{
				"192.168.1.0/24",
			},
		},

		{
			name:  "over one subnet 16 all subnet",
			start: "192.168.0.0",
			end:   "192.168.255.255",
			expected: []string{
				"192.168.0.0/16",
			},
		},

		{
			name:  "over one subnet 8 all subnet",
			start: "192.0.0.0",
			end:   "192.255.255.255",
			expected: []string{
				"192.0.0.0/8",
			},
		},

		{
			name:  "over multiple subnets start from zero",
			start: "192.168.1.0",
			end:   "192.168.3.5",
			expected: []string{
				"192.168.1.0/24",
				"192.168.2.0/24",
				"192.168.3.0/30",
				"192.168.3.4/31",
			},
		},

		{
			name:  "over multiple subnets start from not zero",
			start: "192.168.1.22",
			end:   "192.168.3.5",
			expected: []string{
				"192.168.1.22/31",
				"192.168.1.24/29",
				"192.168.1.32/27",
				"192.168.1.64/26",
				"192.168.1.128/25",
				"192.168.2.0/24",
				"192.168.3.0/30",
				"192.168.3.4/31",
			},
		},

		{
			name:  "over multiple subnets start from not zero two 24 subnets",
			start: "192.168.1.21",
			end:   "192.168.4.5",
			expected: []string{
				"192.168.1.21/32",
				"192.168.1.22/31",
				"192.168.1.24/29",
				"192.168.1.32/27",
				"192.168.1.64/26",
				"192.168.1.128/25",
				"192.168.2.0/23",
				"192.168.4.0/30",
				"192.168.4.4/31",
			},
		},

		{
			name:  "over multiple subnets start from not zero one 16 subnet",
			start: "192.168.1.21",
			end:   "192.170.4.5",
			expected: []string{
				"192.168.1.21/32",
				"192.168.1.22/31",
				"192.168.1.24/29",
				"192.168.1.32/27",
				"192.168.1.64/26",
				"192.168.1.128/25",
				"192.168.2.0/23",
				"192.168.4.0/22",
				"192.168.8.0/21",
				"192.168.16.0/20",
				"192.168.32.0/19",
				"192.168.64.0/18",
				"192.168.128.0/17",
				"192.169.0.0/16",
				"192.170.0.0/22",
				"192.170.4.0/30",
				"192.170.4.4/31",
			},
		},

		{
			name:  "over multiple subnets start from not zero one 8 subnet",
			start: "192.168.1.31",
			end:   "194.170.3.7",
			expected: []string{
				"192.168.1.31/32",
				"192.168.1.32/27",
				"192.168.1.64/26",
				"192.168.1.128/25",
				"192.168.2.0/23",
				"192.168.4.0/22",
				"192.168.8.0/21",
				"192.168.16.0/20",
				"192.168.32.0/19",
				"192.168.64.0/18",
				"192.168.128.0/17",
				"192.169.0.0/16",
				"192.170.0.0/15",
				"192.172.0.0/14",
				"192.176.0.0/12",
				"192.192.0.0/10",
				"193.0.0.0/8",
				"194.0.0.0/9",
				"194.128.0.0/11",
				"194.160.0.0/13",
				"194.168.0.0/15",
				"194.170.0.0/23",
				"194.170.2.0/24",
				"194.170.3.0/29",
			},
		},
	}

	runTestSplitRange(t, cases)
}

func TestSplitRangeV6(t *testing.T) {
	cases := []splitRangeToSubnetsTest{
		{
			name:  "one address",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
			expected: []string{
				"fe80::202:b3ff:fe1e:8329/128",
			},
		},

		{
			name:  "two addresses",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832a",
			expected: []string{
				"fe80::202:b3ff:fe1e:8329/128",
				"fe80::202:b3ff:fe1e:832a/128",
			},
		},

		{
			name:  "three addresses",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832b",
			expected: []string{
				"fe80::202:b3ff:fe1e:8329/128",
				"fe80::202:b3ff:fe1e:832a/127",
			},
		},

		{
			name:  "four addresses",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8329",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832c",
			expected: []string{
				"fe80::202:b3ff:fe1e:8329/128",
				"fe80::202:b3ff:fe1e:832a/127",
				"fe80::202:b3ff:fe1e:832c/128",
			},
		},

		{
			name:  "over one subnet from zero",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8300",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832c",
			expected: []string{
				"fe80::202:b3ff:fe1e:8300/123",
				"fe80::202:b3ff:fe1e:8320/125",
				"fe80::202:b3ff:fe1e:8328/126",
				"fe80::202:b3ff:fe1e:832c/128",
			},
		},

		{
			name:  "over one subnet from one",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8301",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832c",
			expected: []string{
				"fe80::202:b3ff:fe1e:8301/128",
				"fe80::202:b3ff:fe1e:8302/127",
				"fe80::202:b3ff:fe1e:8304/126",
				"fe80::202:b3ff:fe1e:8308/125",
				"fe80::202:b3ff:fe1e:8310/124",
				"fe80::202:b3ff:fe1e:8320/125",
				"fe80::202:b3ff:fe1e:8328/126",
				"fe80::202:b3ff:fe1e:832c/128",
			},
		},

		{
			name:  "over one subnet in middle not align",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8311",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:832c",
			expected: []string{
				"fe80::202:b3ff:fe1e:8311/128",
				"fe80::202:b3ff:fe1e:8312/127",
				"fe80::202:b3ff:fe1e:8314/126",
				"fe80::202:b3ff:fe1e:8318/125",
				"fe80::202:b3ff:fe1e:8320/125",
				"fe80::202:b3ff:fe1e:8328/126",
				"fe80::202:b3ff:fe1e:832c/128",
			},
		},

		{
			name:  "over one subnet in middle with align",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:8310",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:831f",
			expected: []string{
				"fe80::202:b3ff:fe1e:8310/124",
			},
		},

		{
			name:  "over one subnet 112 all subnet",
			start: "fe80:0000:0000:0000:0202:b3ff:fe1e:0000",
			end:   "fe80:0000:0000:0000:0202:b3ff:fe1e:ffff",
			expected: []string{
				"fe80::202:b3ff:fe1e:0/112",
			},
		},

		{
			name:  "over one subnet 112 with zero all subnet",
			start: "fe80:0000:0000:0000:0202:b3ff:0000:0000",
			end:   "fe80:0000:0000:0000:0202:b3ff:0000:ffff",
			expected: []string{
				"fe80::202:b3ff:0:0/112",
			},
		},

		{
			name:  "over one subnet 112 with zeros all subnet",
			start: "fe80:0000:0000:0000:0202:0000:0000:0000",
			end:   "fe80:0000:0000:0000:0202:0000:0000:ffff",
			expected: []string{
				"fe80::202:0:0:0/112",
			},
		},

		{
			name:  "over one subnet 112 with all zeros all subnet",
			start: "fe80:0000:0000:0000:0000:0000:0000:0000",
			end:   "fe80:0000:0000:0000:0000:0000:0000:ffff",
			expected: []string{
				"fe80::/112",
			},
		},

		{
			name:  "over one subnet 96",
			start: "2001:db8:abcd::",
			end:   "2001:db8:abcd::ffff:ffff",
			expected: []string{
				"2001:db8:abcd::/96",
			},
		},

		{
			name:  "over one subnet 64",
			start: "2001:db8:835a:b959::",
			end:   "2001:db8:835a:b959:ffff:ffff:ffff:ffff",
			expected: []string{
				"2001:db8:835a:b959::/64",
			},
		},

		{
			name:  "over one subnet 48",
			start: "2001:db8:1234:0000:0000:0000:0000:0000",
			end:   "2001:db8:1234:ffff:ffff:ffff:ffff:ffff",
			expected: []string{
				"2001:db8:1234::/48",
			},
		},
	}

	runTestSplitRange(t, cases)
}

func runTestSplitRange(t *testing.T, cases []splitRangeToSubnetsTest) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := mustParseAddress(t, c.start)
			e := mustParseAddress(t, c.end)

			iter, err := ip.SplitRangeToSubnets(s, e)
			require.NoError(t, err, "should not returns error for create iter")

			res := make([]string, 0)

			for nt := range iter {
				res = append(res, nt.String())
			}

			require.Equal(t, c.expected, res, "invalid subnets")
		})
	}
}

func mustParseAddress(t *testing.T, s string) netip.Addr {
	resIP, err := netip.ParseAddr(s)
	require.NoError(t, err, "ip '%s' should parse", s)

	return resIP
}
