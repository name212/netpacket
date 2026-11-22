// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"

	"github.com/name212/netpacket/l7/http/client/pool"
	"github.com/name212/netpacket/utils/value"
)

var (
	ErrTimeout           = fmt.Errorf("timeout")
	ErrInvalidStatusCode = fmt.Errorf("invalid status code")
)

type (
	URLPreparator func(*url.URL) *url.URL
	StatusChecker func(code int) bool
)

type Addr interface {
	Addr() string

	IsUnix() bool
	ExtractUnixSocketPath() (string, error)
}

type Client struct {
	transport   http.RoundTripper
	clientsPool *pool.Pool[*http.Client]

	hostURLPreparator URLPreparator
	urlPreparator     URLPreparator

	logger *slog.Logger

	isUnix bool
	host   string

	requestOpts []RequestOpt
}

func NewClient(addr Addr, opts ...ClientOpt) (*Client, error) {
	if value.IsNil(addr) {
		return nil, fmt.Errorf("pass nil Addr to http client")
	}

	passedHost := addr.Addr()

	if passedHost == "" {
		return nil, fmt.Errorf("pass empty host to http client")
	}

	dialHost := passedHost

	clientOpts := newClientOpts()
	for _, o := range opts {
		o(clientOpts)
	}

	logger := clientOpts.logger
	if logger == nil {
		logger = slog.New(logger.Handler())
	}

	hostForURL := addr.Addr()
	dialProtocol := "tcp"
	isUnix := false

	if addr.IsUnix() {
		unixPath, err := addr.ExtractUnixSocketPath()
		if err != nil {
			return nil, err
		}

		dialProtocol = "unix"
		dialHost = unixPath
		hostForURL = "127.0.0.1"
		isUnix = true
	}

	defaultTransport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, fmt.Errorf("http.DefaultTransport is not *http.Transport")
	}

	resultTransport := defaultTransport.Clone()
	defaultDialContext := defaultTransport.DialContext

	resultTransport.DialContext = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return defaultDialContext(ctx, dialProtocol, dialHost)
	}

	resultTransport.Proxy = http.ProxyFromEnvironment

	hostURLPreparator := DummyURLPreparator
	if !clientOpts.noReplaceHost {
		hostURLPreparator = func(u *url.URL) *url.URL {
			u.Host = hostForURL
			return u
		}
	}

	c := &Client{
		transport:         resultTransport,
		hostURLPreparator: hostURLPreparator,
		urlPreparator:     clientOpts.urlPreparator,
		logger:            logger,
		isUnix:            isUnix,
		host:              passedHost,
		requestOpts:       clientsOptsToRequest(opts),
	}

	c.clientsPool = pool.NewPool(func() *http.Client {
		return &http.Client{
			Transport: c.transport,
		}
	})

	return c, nil
}

func (c *Client) Get(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodGet,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) Post(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPost,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) Head(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodHead,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) Delete(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodDelete,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) Put(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPut,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) Patch(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPatch,
		c.pathURL(path),
		opts,
	)
}

func (c *Client) pathURL(path string) *url.URL {
	return &url.URL{
		Path: path,
	}
}

func (c *Client) getRequest(method string, u *url.URL, opts []RequestOpt) (*http.Request, *RequestOpts, error) {
	if u == nil {
		return nil, nil, fmt.Errorf("empty url for host '%s' with method %s", c.host, method)
	}

	if u.Path == "" {
		return nil, nil, fmt.Errorf("empty url for host '%s' with method %s", c.host, method)
	}

	passedPath := u.Path

	requestOpts := newRequestOpts()

	for _, o := range c.requestOpts {
		o(requestOpts)
	}

	for _, o := range opts {
		o(requestOpts)
	}

	rURL := copyURL(u)
	rURL.Scheme = requestOpts.schema
	rURL = c.hostURLPreparator(rURL)
	rURL = c.urlPreparator(rURL)

	if _, err := url.Parse(rURL.String()); err != nil {
		return nil, nil, c.err("url is incorrect", err, method, passedPath)
	}

	res, err := http.NewRequest(method, u.String(), requestOpts.body)
	if err != nil {
		return nil, nil, c.err("cannot create http request", err, method, passedPath)
	}

	if len(requestOpts.headers) > 0 {
		for k, v := range requestOpts.headers {
			res.Header.Add(k, v)
		}
	}

	return res, requestOpts, nil
}

func (c *Client) do(parent context.Context, method string, u *url.URL, opts []RequestOpt) (*ResponseBody, error) {
	r, requestOpts, err := c.getRequest(method, u, opts)
	if err != nil {
		return nil, err
	}

	cl := c.clientsPool.Get()
	if cl == nil {
		return nil, c.err("cannot get client", fmt.Errorf("pool return nil"), method, u.Path)
	}

	defer func() {
		c.clientsPool.Put(cl)
	}()

	c.log(requestOpts, "start request", logFuncInfo, r, 0)

	ctx, cancel := context.WithTimeout(parent, requestOpts.timeout)
	defer cancel()

	r = r.WithContext(ctx)

	response, err := cl.Do(r)

	if err != nil {
		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return nil, c.errRequest(requestOpts, FailRequestMsg, ErrTimeout, r, 0)
			}
		}

		return nil, c.errRequest(requestOpts, FailRequestMsg, err, r, 0)
	}

	var rBody *ResponseBody
	if !value.IsNil(response.Body) {
		rBody = newResponseBody(response.Body, r.URL.Host, method, requestOpts.autoCloseResponseBody)
	}

	if !requestOpts.statusChecker(response.StatusCode) {
		return rBody, c.errRequest(requestOpts, ErrInvalidStatusCode.Error(), ErrInvalidStatusCode, r, response.StatusCode)
	}

	c.log(requestOpts, OKRequestMsg, logFuncInfo, r, response.StatusCode)

	return rBody, nil
}

func (c *Client) err(msg string, errR error, method, passedPath string) error {
	return fmt.Errorf(
		"%s '%s %s' to server '%s': %w",
		msg,
		method,
		passedPath,
		c.host,
		errR,
	)
}

func (c *Client) errRequest(requestOpts *RequestOpts, msg string, errR error, r *http.Request, code int) error {
	res := c.err(msg, errR, r.Method, r.URL.Path)
	if code > 0 {
		res = fmt.Errorf("%w response code %d", res, code)
	}

	c.log(requestOpts, msg, logFuncErr, r, code, slog.String(ErrAttr, errR.Error()))

	return res
}

type logFunc func(string, ...any)
type logFuncConsumer func(l *slog.Logger) logFunc

func logFuncErr(l *slog.Logger) logFunc {
	return l.Error
}

func logFuncInfo(l *slog.Logger) logFunc {
	return l.Info
}

func (c *Client) log(requestOpts *RequestOpts, msg string, fc logFuncConsumer, r *http.Request, code int, attrs ...any) {
	logger := requestOpts.logger
	if logger == nil {
		logger = c.logger
	}

	if logger == nil {
		return
	}

	a := []any{
		slog.String(LogTypeAttr, LogTypeClientVal),
		slog.String(MethodAttr, r.Method),
		slog.String(PathAttr, r.URL.Path),
		slog.String(HostAttr, c.host),
	}

	if id := requestOpts.getUserID(); id != "" {
		a = append(a, slog.String(UserAttr, id))
	}

	if agent := requestOpts.getUserAgent(); agent != "" {
		a = append(a, slog.String(UserAgentAttr, agent))
	}

	if code > 0 {
		a = append(a, slog.Int(StatusCodeAttr, code))
	}

	if c.isUnix {
		a = append(a, slog.Bool(IsUnixSocketAttr, true))
	}

	a = append(a, attrs...)

	fc(logger)(msg, a...)
}
