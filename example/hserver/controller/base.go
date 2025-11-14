package controller

import (
	"encoding/json"
	"fmt"
	"github.com/tonly18/xws/core/logger"
	"github.com/tonly18/xws/core/xerror"
	"github.com/tonly18/xws/core/ziface"
	"github.com/tonly18/xws/core/zserver"
	"github.com/tonly18/xws/example/global"
	"net/http"
	"runtime"

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
		//connId := r.Header.Get("conn_id")
		userId := r.Header.Get("user_id")
		clientIP := r.Header.Get("client_ip")
		traceId := r.Header.Get("trace_id")

		//参数判断
		if clientIP == "" || traceId == "" {
			writeResponseData(w, &zserver.Response{Code: "1000"})
			return
		}

		//获取当前玩家conn
		var conn ziface.IConnection
		uid := cast.ToInt64(userId)
		if uid > 0 {
			connection, err := global.GetWsServer().GetConnMgr().GetByUid(uid)
			if err == nil && connection != nil {
				conn = connection
			}
		}

		//request
		request := &zserver.Request{
			ResponseWriter: w,
			Request:        r,
			UserID:         uid,
			Conn:           conn,
		}
		request.SetData("user_id", userId)
		request.SetData("client_ip", clientIP)
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
