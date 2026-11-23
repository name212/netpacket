// Copyright 2026
// license that can be found in the LICENSE file.

package client

import "net/url"

func CopyURL(u *url.URL) *url.URL {
	n, err := url.Parse(u.String())
	if err != nil {
		n = copyURL(u)
	}

	return n
}

func DummyURLPreparator(u *url.URL) *url.URL {
	return u
}

func copyURL(u *url.URL) *url.URL {
	return &url.URL{
		Scheme:      u.Scheme,
		Host:        u.Host,
		User:        u.User,
		Path:        u.Path,
		RawPath:     u.RawPath,
		RawQuery:    u.RawQuery,
		Opaque:      u.Opaque,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
		OmitHost:    u.OmitHost,
		ForceQuery:  u.ForceQuery,
	}
}
