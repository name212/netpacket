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
	"github.com/name212/netpacket/l7/http/content"
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
	if logger != nil {
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
	} else {
		host, port, err := net.SplitHostPort(dialHost)
		if err != nil {
			return nil, fmt.Errorf("cannot extract host from '%s': %w", dialHost, err)
		}

		switch host {
		case "0.0.0.0":
			hst := net.JoinHostPort("127.0.0.1", port)
			dialHost = hst
			hostForURL = hst
			passedHost = hst
		case "::":
			hst := net.JoinHostPort("::", port)
			dialHost = hst
			hostForURL = hst
			passedHost = hst
		}
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

func NewClientForURL(u *url.URL, opts ...ClientOpt) (*Client, error) {
	addr, err := NewURLAddr(u)
	if err != nil {
		return nil, fmt.Errorf("cannot create client for url: %w", err)
	}

	passOpts := opts

	if u.Scheme == "https" {
		passOpts = append(passOpts, WithClientHTTPS())
	}

	return NewClient(addr, passOpts...)
}

func (c *Client) Get(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodGet,
		c.pathURL(path),
		opts,
	)
}

func Get(parent context.Context, addr Addr, path string, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return c.Get(parent, path, opts...)
}

func GetURL(parent context.Context, u *url.URL, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClientForURL(u)
	if err != nil {
		return nil, err
	}

	return c.Get(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) Post(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPost,
		c.pathURL(path),
		opts,
	)
}

func Post(parent context.Context, addr Addr, path string, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return c.Post(parent, path, opts...)
}

func PostURL(parent context.Context, u *url.URL, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClientForURL(u)
	if err != nil {
		return nil, err
	}

	return c.Post(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) Head(parent context.Context, path string, opts ...RequestOpt) error {
	r, err := c.do(
		parent,
		http.MethodHead,
		c.pathURL(path),
		opts,
	)

	err = appendErr(err, r.Close())

	return err
}

func Head(parent context.Context, addr Addr, path string, opts ...RequestOpt) error {
	c, err := NewClient(addr)
	if err != nil {
		return err
	}

	return c.Head(parent, path, opts...)
}

func HeadURL(parent context.Context, u *url.URL, opts ...RequestOpt) error {
	c, err := NewClientForURL(u)
	if err != nil {
		return err
	}

	return c.Head(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) Delete(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodDelete,
		c.pathURL(path),
		opts,
	)
}

func Delete(parent context.Context, addr Addr, path string, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return c.Delete(parent, path, opts...)
}

func DeleteURL(parent context.Context, u *url.URL, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClientForURL(u)
	if err != nil {
		return nil, err
	}

	return c.Delete(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) Put(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPut,
		c.pathURL(path),
		opts,
	)
}

func Put(parent context.Context, addr Addr, path string, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return c.Put(parent, path, opts...)
}

func PutURL(parent context.Context, u *url.URL, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClientForURL(u)
	if err != nil {
		return nil, err
	}

	return c.Put(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) Patch(parent context.Context, path string, opts ...RequestOpt) (*ResponseBody, error) {
	return c.do(
		parent,
		http.MethodPatch,
		c.pathURL(path),
		opts,
	)
}

func Patch(parent context.Context, addr Addr, path string, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClient(addr)
	if err != nil {
		return nil, err
	}

	return c.Patch(parent, path, opts...)
}

func PatchURL(parent context.Context, u *url.URL, opts ...RequestOpt) (*ResponseBody, error) {
	c, err := NewClientForURL(u)
	if err != nil {
		return nil, err
	}

	return c.Patch(parent, GetFullPathFromURL(u), opts...)
}

func (c *Client) pathURL(path string) *url.URL {
	return &url.URL{
		Path: path,
	}
}

func (c *Client) prepareRequestOpts(method string, u *url.URL, opts []RequestOpt) (*RequestOpts, string, error) {
	if u == nil {
		return nil, "", fmt.Errorf("empty url for host '%s' with method %s", c.host, method)
	}

	if u.Path == "" {
		return nil, "", fmt.Errorf("empty url for host '%s' with method %s", c.host, method)
	}

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

	urlStr := rURL.String()

	if _, err := url.Parse(urlStr); err != nil {
		return nil, "", c.err("url is incorrect", err, method, u.Path)
	}

	return requestOpts, urlStr, nil
}

func (c *Client) getRequest(method, rURL string, requestOpts *RequestOpts) (*http.Request, error) {
	res, err := http.NewRequest(method, rURL, requestOpts.body)
	if err != nil {
		return nil, err
	}

	if len(requestOpts.headers) > 0 {
		for k, v := range requestOpts.headers {
			res.Header.Add(k, v)
		}
	}

	return res, nil
}

func (c *Client) do(parent context.Context, method string, u *url.URL, opts []RequestOpt) (*ResponseBody, error) {
	requestOpts, urlStr, err := c.prepareRequestOpts(method, u, opts)
	if err != nil {
		return nil, err
	}

	passedPath := u.Path
	if passedPath == "" {
		passedPath = "/"
	}

	cl := c.clientsPool.Get()
	if cl == nil {
		return nil, c.err("cannot get client", fmt.Errorf("pool return nil"), method, u.Path)
	}

	defer func() {
		c.clientsPool.Put(cl)
	}()

	var respBody *ResponseBody

	c.setRetrierLogger(requestOpts, method, passedPath)

	for {
		r, err := c.getRequest(method, urlStr, requestOpts)
		if err != nil {
			return nil, c.err("cannot create http request", err, method, passedPath)
		}

		response, err := c.runRequest(parent, cl, r, requestOpts)

		if err != nil {
			var retryErr error
			if requestOpts.retrier != nil {
				retryErr = requestOpts.retrier.makeRetry(parent, err, 0)
				if retryErr == nil {
					continue
				}
			}

			if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
				return nil, c.errRequest(requestOpts, LogFailRequestMsg, appendErr(ErrTimeout, retryErr), r, 0)
			}

			return nil, c.errRequest(requestOpts, LogFailRequestMsg, appendErr(err, retryErr), r, 0)
		}

		statusCode := response.StatusCode

		if !value.IsNil(response.Body) {
			respBody = newResponseBody(
				response.Body,
				r.URL.Host,
				method,
				requestOpts.autoCloseResponseBody,
				content.GetResponseContentType(response),
			)
		}

		var retryErr error

		validStatus := requestOpts.statusChecker(statusCode)
		if !validStatus {
			if requestOpts.retrier != nil {
				retryErr = requestOpts.retrier.makeRetry(parent, err, statusCode)
				if retryErr == nil {
					_ = c.handleResponseErr(requestOpts, respBody, r)
					continue
				}
			}

			respBody.withLogger(c.readBodyLogger(requestOpts, r, statusCode))
		} else {
			respBody.withLogger(c.readBodyLogger(requestOpts, r, statusCode))
		}

		c.writeResponseHeaders(response, requestOpts)

		if !validStatus {
			msg := fmt.Sprintf("%v: %d", ErrInvalidStatusCode, statusCode)

			err := ErrInvalidStatusCode

			if retryErr != nil && retryErr == ErrInvalidStatusCode {
				err = nil
			}

			err = appendErr(err, retryErr)

			err = appendErr(err, c.handleResponseErr(requestOpts, respBody, r))

			return respBody, c.errRequest(requestOpts, msg, err, r, statusCode)
		}

		c.log(requestOpts, LogOKRequestMsg, LogFuncInfo, r, statusCode)
		break
	}

	return respBody, nil
}

func (c *Client) writeResponseHeaders(r *http.Response, opts *RequestOpts) {
	if opts.responseHeadersDest == nil {
		return
	}

	if len(r.Header) == 0 {
		return
	}

	for k := range r.Header {
		if !opts.responseHeadersFilter(k) {
			continue
		}

		opts.responseHeadersDest[k] = r.Header.Get(k)
	}
}

func (c *Client) handleResponseErr(requestOpts *RequestOpts, respBody *ResponseBody, r *http.Request) error {
	if respBody == nil {
		return nil
	}

	c.logBody(requestOpts, respBody, r)

	if requestOpts.autoCloseResponseBody {
		if err := respBody.Close(); err != nil {
			c.log(requestOpts, "cannot close response body", LogFuncErr, r, 0)
			return err
		}
	}

	return nil
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

	c.log(requestOpts, msg, LogFuncErr, r, code, slog.String(LogErrAttr, errR.Error()))

	return res
}

func (c *Client) runRequest(parent context.Context, cl *http.Client, r *http.Request, requestOpts *RequestOpts) (*http.Response, error) {
	c.log(requestOpts, "start request", LogFuncInfo, r, 0)

	ctx, cancel := context.WithTimeout(parent, requestOpts.timeout)
	defer cancel()

	r = r.WithContext(ctx)

	return cl.Do(r)
}

func appendErr(err, err2 error) error {
	if err2 == nil {
		return err
	}

	if err == nil {
		return nil
	}

	return fmt.Errorf("%w; %w", err, err2)
}
