// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"fmt"
	"io"
)

var (
	ErrClosed = fmt.Errorf("response body already closed")
)

type ResponseBody struct {
	reader    io.ReadCloser
	url       string
	method    string
	needClose bool

	closed bool
}

func newResponseBody(r io.ReadCloser, url, method string, needClose bool) *ResponseBody {
	return &ResponseBody{
		reader:    r,
		url:       url,
		method:    method,
		needClose: needClose,
		closed:    false,
	}
}

func (r *ResponseBody) Read(p []byte) (int, error) {
	if r == nil {
		return 0, nil
	}

	const errMsg = "as reader"

	if err := r.isClosed(errMsg); err != nil {
		return 0, err
	}

	n, err := r.reader.Read(p)
	if err != nil {
		return n, r.readError(errMsg, err, n)
	}

	return n, nil
}

func (r *ResponseBody) ReadBytes() ([]byte, error) {
	if r == nil {
		return nil, nil
	}

	return r.readBytes("bytes")
}

func (r *ResponseBody) ReadString() (string, error) {
	if r == nil {
		return "", nil
	}

	res, err := r.readBytes("string")
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (r *ResponseBody) Copy(w io.Writer) (int64, error) {
	if r == nil {
		return 0, nil
	}

	const errMsg = "copy to dest"

	if err := r.isClosed(errMsg); err != nil {
		return 0, err
	}

	n, err := io.Copy(w, r.reader)
	if err != nil {
		err = r.readError(errMsg, err)

		if r.needClose {
			err = r.appendErrClose(err, r.close())
		}

		return n, err
	}

	if r.needClose {
		err = r.appendErrClose(nil, r.close())
	}

	return n, err
}

func (r *ResponseBody) Close() error {
	if r == nil {
		return nil
	}

	if r.closed {
		return nil
	}

	defer func() {
		r.closed = true
	}()

	if err := r.close(); err != nil {
		r.closed = true

		return r.appendErrClose(nil, err)
	}

	r.closed = true

	return nil
}

func (r *ResponseBody) IsClosed() error {
	if r == nil {
		return nil
	}

	return r.isClosed("")
}

func (r *ResponseBody) readBytes(msg string) ([]byte, error) {
	if err := r.isClosed(msg); err != nil {
		return nil, err
	}

	res, err := io.ReadAll(r.reader)
	if err != nil {
		err = r.readError(msg, err)

		if r.needClose {
			err = r.appendErrClose(err, r.close())
		}

		return nil, err
	}

	if r.needClose {
		err = r.appendErrClose(nil, r.close())
	}

	return res, err
}

func (r *ResponseBody) close() error {
	return r.reader.Close()
}

func (r *ResponseBody) isClosed(msg string) error {
	if r.closed {
		if msg != "" {
			return r.readError(msg, ErrClosed)
		}

		return fmt.Errorf(
			"%W for '%s %s'",
			ErrClosed,
			r.method,
			r.url,
		)
	}

	return nil
}

func (r *ResponseBody) readError(msg string, err error, n ...int) error {
	btMsg := ""
	if len(n) > 0 {
		btMsg = fmt.Sprintf(" read bytes %d", n[0])
	}

	if msg != "" {
		msg = " " + msg
	}

	return fmt.Errorf(
		"cannot read body%s for '%s %s'%s: %w",
		msg,
		r.method,
		r.url,
		btMsg,
		err,
	)
}

func (r *ResponseBody) appendErrClose(err, errClose error) error {
	var curErr error

	if errClose != nil {
		curErr = fmt.Errorf("cannot close body for '%s %s': %w", r.method, r.url, errClose)
	}

	if err != nil {
		if curErr == nil {
			return err
		}

		curErr = fmt.Errorf("%v\n%v", err, curErr)
	}

	return curErr
}
