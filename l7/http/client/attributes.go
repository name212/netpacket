// Copyright 2026
// license that can be found in the LICENSE file.

package client

import "github.com/name212/netpacket/l7/http/log"

const (
	HostAttr         = "host"
	IsUnixSocketAttr = "unix_socket"
	MethodAttr       = log.MethodAttr
	PathAttr         = log.PathAttr
	UserAttr         = log.UserAttr
	UserAgentAttr    = log.UserAgentAttr
	RequestLenAttr   = log.RequestLenAttr
	ResponseLenAttr  = log.ResponseLenAttr
	StatusCodeAttr   = log.StatusCodeAttr
	LogTypeAttr      = log.LogTypeAttr

	ErrAttr = log.ErrAttr
)

const (
	LogTypeClientVal = "http-client"

	OKRequestMsg   = "Request OK"
	FailRequestMsg = "Request error"
)
