package controller

import (
	"encoding/json"
	"fmt"

	"net/http"
	"runtime"

	"github.com/extra-time-zone/xws/core/logger"
	"github.com/extra-time-zone/xws/core/xerror"
	"github.com/extra-time-zone/xws/core/ziface"
	"github.com/extra-time-zone/xws/core/zserver"
	"github.com/google/uuid"

	"github.com/spf13/cast"
)

func WrapHandle(handler func(*zserver.Request) (*zserver.Response, xerror.Error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(r.Context(), fmt.Sprintf(`[wrap handle] Error(1): %v`, err))
				logger.Error(r.Context(), fmt.Sprintf(`[wrap handle] ProxyId:%v, ServerId:%v, UserId: %v, ClientIP: %v`, r.Header.Get("proxy_id"), r.Header.Get("server_id"), r.Header.Get("user_id"), r.Header.Get("client_ip")))
				for i := 1; i < 20; i++ {
					if pc, file, line, ok := runtime.Caller(i); ok {
						fcName := runtime.FuncForPC(pc).Name()
						logger.Error(r.Context(), fmt.Sprintf(`[wrap handle] goroutine:%v, file:%s, function name:%s, line:%d`, pc, file, fcName, line))
					}
				}
				logger.Error(r.Context(), fmt.Sprintf(`[wrap handle] Error(2): %v`, err))
			}
		}()

		//params
		var userId int64
		var conn []ziface.IConnection
		var traceId = r.Header.Get("trace_id")
		if uid := r.Header.Get("user_id"); uid != "" {
			userId = cast.ToInt64(uid)
		}
		if traceId == "" {
			traceId = uuid.NewString()
		}

		//request
		request := &zserver.Request{
			ResponseWriter: w,
			Request:        r,
			UserID:         userId,
			Conn:           conn,
		}
		request.SetData("user_id", userId)
		//request.SetData("client_ip", clientIP)
		request.SetData("trace_id", traceId)

		//handler
		resp, xerr := handler(request)
		if xerr != nil {
			logger.Error(request, fmt.Sprintf(`[code:%v, data:%v, message:%v]`, resp.Code, resp.Data, resp.Message))
		} else {
			logger.Info(request, fmt.Sprintf(`[code:%v, data:%v, message:%v]`, resp.Code, resp.Data, resp.Message))
		}

		//result
		writeResponseData(w, resp)
	}
}

func writeResponseData(w http.ResponseWriter, params *zserver.Response) {
	dataByte, _ := json.Marshal(params)
	w.Header().Set("content-length", cast.ToString(len(dataByte)))
	w.Write(dataByte)
	w.(http.Flusher).Flush()
}

func writeResponseBytes(w http.ResponseWriter, data []byte) {
	w.Header().Set("content-length", cast.ToString(len(data)))
	w.Write(data)
	w.(http.Flusher).Flush()
}
