// Copyright 2026
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"

	"github.com/name212/netpacket/l7/http/content"
)

const (
	UserIDHeader      = "X-User-ID"
	NotFoundUserIDMsg = "cannot extract user id"
)

type UserIDConsumer func(*http.Request) (string, error)

func SetCtxUserID(next http.Handler, consumer ...UserIDConsumer) http.Handler {
	c := getUserIDConsumer(consumer)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r = setUserID(w, r, c); r != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func getUserIDConsumer(consumer []UserIDConsumer) UserIDConsumer {
	c := ReadUserNameFromHeader
	if len(consumer) > 0 && consumer[0] != nil {
		c = consumer[0]
	}

	return c
}

func setUserID(w http.ResponseWriter, r *http.Request, c UserIDConsumer) *http.Request {
	data, r := getAndSetServerData(r)

	id, err := c(r)
	if err != nil {
		logger := data.Logger()

		logger.RequestError(r, NotFoundUserIDMsg, err)

		out, ct := data.ErrorWriter()(NotFoundUserIDMsg, err)
		w.Header().Add(content.ContentTypeHeader, ct)

		w.WriteHeader(http.StatusForbidden)

		if _, errW := w.Write(out); errW != nil {
			logger.RequestError(r, "cannot write not found user id response", nil)
		}

		return nil
	}

	data.userID = id
	return r
}

func ReadUserNameFromHeader(r *http.Request) (string, error) {
	res := r.Header.Get(UserIDHeader)
	if res == "" {
		res = UnknownVal
	}

	return res, nil
}
