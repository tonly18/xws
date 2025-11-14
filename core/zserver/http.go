package zserver

import (
	"context"
	"fmt"
	"github.com/tonly18/xws/core/logger"
	"net/http"
	"time"
)

type httpServer struct {
	Http   http.Server
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewHttpServer(conf *HttpServerConfig) *httpServer {
	return &httpServer{
		Http: http.Server{
			Addr:           fmt.Sprintf(`%s:%d`, conf.IP, conf.Port),
			Handler:        conf.Handler,
			ReadTimeout:    5 * time.Second,  //从链接被接受开始,到request body完全读取完为止.
			WriteTimeout:   10 * time.Second, //http:从request head读取结束开始到response write完成为止.
			IdleTimeout:    10 * time.Second, //空闲时长:如果IdleTimeout为0,则使用ReadTimeout.
			MaxHeaderBytes: 1 << 20,
		},
	}
}

func (s *httpServer) Start() error {
	logger.LogInfof("[HTTP SERVER START] Successful. Listening at Addr: %s", s.Http.Addr)

	return s.Http.ListenAndServe()
}

func (s *httpServer) Stop() {
	s.Ctx, s.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer s.Cancel()

	_ = s.Http.Shutdown(s.Ctx)
}
