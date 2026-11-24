// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"net/url"

	"github.com/name212/netpacket/l7/params/unix"
)

var _ Addr = &URLAddr{}

type URLAddr struct {
	url *url.URL
}

func NewURLAddr(u *url.URL) (*URLAddr, error) {
	if u == nil {
		return nil, fmt.Errorf("passed nil url to NewURLAddr")
	}

	if u.Host == "" {
		return nil, fmt.Errorf("host url is empty for NewURLAddr")
	}

	cpy, err := url.Parse(u.String())
	if err != nil {
		return nil, fmt.Errorf("invalid url for NewURLAddr: %w", err)
	}

	return &URLAddr{
		url: cpy,
	}, nil
}

func (a *URLAddr) URL() *url.URL {
	return a.url
}

func (a *URLAddr) Addr() string {
	return a.url.Host
}

func (a *URLAddr) ExtractUnixSocketPath() (string, error) {
	return unix.ExtractUnixSocketPath(a.Addr())
}

func (a *URLAddr) IsUnix() bool {
	return unix.IsUnix(a.Addr())
}
