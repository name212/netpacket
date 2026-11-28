// Copyright 2026
// license that can be found in the LICENSE file.

package client

import (
	"encoding/json"
	"fmt"

	"github.com/name212/netpacket/l7/http/content"
)

var (
	ErrUnmarshalerNotFound = fmt.Errorf("unmarshaler not found")

	defaultUnmarshalerConsumer = NewUnmarshalerConsumer(json.Unmarshal)
)

type Unmarshaler func([]byte, any) error

type UnmarshalerConsumer struct {
	m map[content.ContentType]Unmarshaler
}

func SetDefaultUnmarshalerConsumer(c *UnmarshalerConsumer) {
	if c != nil {
		defaultUnmarshalerConsumer = c
	}
}

func NewUnmarshalerConsumer(textUnmarshaler ...Unmarshaler) *UnmarshalerConsumer {
	r := &UnmarshalerConsumer{
		m: make(map[content.ContentType]Unmarshaler),
	}

	var txt Unmarshaler

	if len(textUnmarshaler) > 0 && textUnmarshaler[0] != nil {
		txt = textUnmarshaler[0]
	}

	r.Add(content.TypeText, txt)
	r.Add(content.TypeNo, txt)
	r.Add(content.TypeJSON, json.Unmarshal)

	return r
}

func (c *UnmarshalerConsumer) Add(ct content.ContentType, f Unmarshaler) *UnmarshalerConsumer {
	if f != nil {
		c.m[ct] = f
	}

	return c
}

func (c *UnmarshalerConsumer) get(ct content.ContentType) (Unmarshaler, error) {
	f, ok := c.m[ct]
	if !ok {
		return nil, fmt.Errorf("content type '%s' %w", ct, ErrUnmarshalerNotFound)
	}
	return f, nil
}
