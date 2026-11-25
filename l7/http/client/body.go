// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	"github.com/name212/netpacket/l7/http/content"
)

var (
	ErrClosed    = fmt.Errorf("response body already closed")
	ErrEmptyBody = fmt.Errorf("response body is empty")
)

type ResponseBody struct {
	reader      io.ReadCloser
	url         string
	method      string
	needClose   bool
	contentType content.ContentType
	logger      *responseBodyReadLogger

	closed bool
}

func newResponseBody(r io.ReadCloser, url, method string, needClose bool, contentType content.ContentType) *ResponseBody {
	return &ResponseBody{
		reader:      r,
		url:         url,
		method:      method,
		needClose:   needClose,
		closed:      false,
		contentType: contentType,
		logger:      nil,
	}
}

//nolint:unparam
func (r *ResponseBody) withLogger(l *responseBodyReadLogger) *ResponseBody {
	r.logger = l
	return r
}

func (r *ResponseBody) GetContentType() content.ContentType {
	return r.contentType
}

func (r *ResponseBody) Read(p []byte) (int, error) {
	if r == nil {
		return 0, nil
	}

	const errMsg = "as bytes buffer reader"

	if err := r.isClosed(errMsg); err != nil {
		return 0, err
	}

	n, err := r.reader.Read(p)

	r.logger.read(int64(n), errMsg)

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

func (r *ResponseBody) ReadAs[T any](c ...*UnmarshalerConsumer) (T, error) {
	var zero T

	msg := fmt.Sprintf("as %s", reflect.TypeFor[T]().String())

	if r == nil {
		return zero, ErrEmptyBody
	}

	fConsumer := defaultUnmarshalerConsumer
	if len(c) > 0 && c[0] != nil {
		fConsumer = c[0]
	}

	f, err := fConsumer.get(r.contentType)
	if err != nil {
		return zero, r.readError(msg, err)
	}

	bt, err := r.readBytes(msg)
	if err != nil {
		return zero, err
	}

	btLen := len(bt)

	if btLen == 0 {
		return zero, r.readError(ErrEmptyBody.Error(), ErrEmptyBody, 0)
	}

	r.logger.read(int64(btLen), msg)

	var res T
	if err := f(bt, &res); err != nil {
		const uMsg = "cannot unmarshal response"
		r.logger.body(bt, fmt.Sprintf("%s %s", uMsg, msg))
		return zero, r.readError(uMsg, err)
	}

	return res, nil
}

func (r *ResponseBody) CopyKeep(w io.Writer, limit ...int64) (int64, error) {
	if r == nil {
		return 0, nil
	}

	const errMsg = "copy keep to dest"

	if err := r.isClosed(errMsg); err != nil {
		return 0, err
	}

	curReader := r.reader

	buf := bytes.NewBuffer(nil)
	r.reader = newReaderClose(buf)

	destWriter := newLimitWriter(w, limit)

	writer := io.MultiWriter(buf, destWriter)

	n, err := io.Copy(writer, curReader)
	n = destWriter.getWritten(n)

	err = r.appendErrClose(err, curReader.Close())
	if err != nil {
		err = r.readError(errMsg, err)
	}

	return n, err
}

func (r *ResponseBody) Copy(w io.Writer, limit ...int64) (int64, error) {
	if r == nil {
		return 0, nil
	}

	const errMsg = "copy to dest"

	if err := r.isClosed(errMsg); err != nil {
		return 0, err
	}

	destWriter := newLimitWriter(w, limit)

	n, err := io.Copy(destWriter, r.reader)
	n = destWriter.getWritten(n)

	r.logger.read(n, errMsg)

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

	r.logger.read(int64(len(res)), msg)

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
