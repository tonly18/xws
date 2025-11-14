package zserver

import (
	"net/http"
	"time"

	"github.com/tonly18/xws/core/ziface"
	"github.com/tonly18/xws/core/zutils"

	"github.com/spf13/cast"
)

// Request 请求
type Request struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	UserID         int64
	Conn           ziface.IConnection
	data           map[string]any
}

func (r *Request) GetData(key string) any {
	return r.data[key]
}

func (r *Request) SetData(key string, value any) {
	if len(r.data) == 0 {
		r.data = make(map[string]any, 10)
	}
	r.data[key] = value
}

func (r *Request) Deadline() (deadline time.Time, ok bool) {
	return r.Request.Context().Deadline()
}

func (r *Request) Done() <-chan struct{} {
	return r.Request.Context().Done()
}

func (r *Request) Err() error {
	return r.Request.Context().Err()
}

func (r *Request) Value(key any) any {
	value := r.GetData(cast.ToString(key))
	if zutils.IsNil(value) {
		value = r.Request.Context().Value(key)
		if zutils.IsNil(value) {
			if r.Conn != nil {
				value = r.Conn.GetProperty(cast.ToString(key))
			}
		}
	}

	return value
}
