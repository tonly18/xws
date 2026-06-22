package znet

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/extra-time-zone/xws/core/logger"
	"github.com/extra-time-zone/xws/core/zconf"
	"github.com/extra-time-zone/xws/core/ziface"
	"github.com/extra-time-zone/xws/core/zutils"

	"github.com/gorilla/websocket"
)

type Server struct {
	env string
	//服务器ID
	id uint32
	//服务器的名称
	name string
	//服务绑定的IP地址
	ip string
	//服务绑定的端口
	port int
	//path
	path string
	//当前Server的消息管理模块，用来绑定MsgID和对应的处理方法
	msgHandler ziface.IMsgHandle
	//当前Server的链接管理器
	connMgr ziface.IConnManager
	//该Server的连接创建时Hook函数
	onConnStart func(conn ziface.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn ziface.IConnection)
	// websocket
	upgrader *websocket.Upgrader

	// ws安全
	safeIPLimiter ziface.IPLimiter

	ctx    context.Context
	cancel context.CancelFunc
	packet ziface.Packet
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		env:        zconf.Config.ENV,
		id:         zconf.Config.ServerID,
		name:       zconf.Config.Name,
		ip:         zconf.Config.Host,
		port:       zconf.Config.Port,
		path:       zconf.Config.Path,
		msgHandler: NewMsgHandle(),
		connMgr:    NewConnManager(),
		packet:     NewPacket(),
		ctx:        nil,
		cancel:     nil,
		upgrader: &websocket.Upgrader{
			EnableCompression: true,
			ReadBufferSize:    4096,
			WriteBufferSize:   4096,
			HandshakeTimeout:  time.Second * 5,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		safeIPLimiter: NewIPLimiter(),
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())

	return s
}

// StartConn Start Conn
func (s *Server) StartConn(conn ziface.IConnection) {
	conn.Start()
}

// Start 开启网络服务
func (s *Server) Start() {
	//0 启动worker工作池机制
	s.msgHandler.StartWorkerPool()

	//开启一个go去做服务端Listener
	go s.ListenWebsocketConn()
}

func (s *Server) ListenWebsocketConn() {
	logger.LogInfof("[WebSocket Server] Successful Server Name:%s, Listen At:%s:%v, Path:%v", s.name, s.ip, s.port, s.path)

	http.HandleFunc(s.path, func(w http.ResponseWriter, r *http.Request) {
		//0. 服务退出,拒绝连接
		if zconf.ServiceExit.Load() {
			logger.LogInfo("[WebSocket Server] ListenWebsocketConn Refuse New Connection")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//1. 单个IP限制连接数量
		clientIp := zutils.GetClientIP(r)
		if !s.safeIPLimiter.Allow(clientIp) {
			logger.LogInfof("[WebSocket Server] ListenWebsocketConn client not allow ip:%v Connection Limited", clientIp)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		//2. 设置服务器最大连接控制,如果超过最大连接,则等待
		if maxConn := s.GetConnMgr().Count(); maxConn >= zconf.Config.MaxConn {
			logger.LogErrorf("[WebSocket Server] ListenWebsocketConn conn limit exceeded maxConn:%d, maxConnConfig:%d", maxConn, zconf.Config.MaxConn)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		//3. 判断 header 里面是有子协议
		if len(r.Header.Get("Sec-Websocket-Protocol")) > 0 {
			s.upgrader.Subprotocols = websocket.Subprotocols(r)
		}

		//4. 升级成 websocket 连接
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.LogErrorf("[WebSocket Server] ListenWebsocketConn Upgrade Websocket fail error:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//5. 处理该新连接请求的业务方法,此时应该有 handler 和 conn是绑定的
		wsConn := newWebsocketConn(s, conn, clientIp, r)

		go s.StartConn(wsConn)
	})

	if zconf.Config.CertFile != "" && zconf.Config.PrivateKeyFile != "" {
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", s.ip, s.port), zconf.Config.CertFile, zconf.Config.PrivateKeyFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.ip, s.port), nil); err != nil {
			panic(err)
		}
	}
}

// Stop 停止服务
func (s *Server) Stop() {
	logger.LogInfo("[WebSocket Server] Server Name:", s.name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.connMgr.Clear()

	//退出
	s.cancel()
}

// Serve 运行服务
func (s *Server) Serve() {
	s.Start()

	//阻塞,否则主Go退出,listener的go将会退出
	select {
	case <-s.ctx.Done():
		logger.LogInfo("[WebSocket Server] Context Cancel")
	}
}

// AddRouter 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.msgHandler.AddRouter(msgID, router)
}

// GetConnMgr 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.connMgr
}

// SetOnConnStart 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.onConnStart = hookFunc
}
func (s *Server) GetOnConnStart() func(ziface.IConnection) {
	return s.onConnStart
}

// SetOnConnStop 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.onConnStop = hookFunc
}
func (s *Server) GetOnConnStop() func(ziface.IConnection) {
	return s.onConnStop
}

func (s *Server) Packet() ziface.Packet {
	return s.packet
}

func (s *Server) GetID() uint32 {
	return s.id
}

// GetMsgHandler 获取Server绑定的消息处理模块
func (s *Server) GetMsgHandler() ziface.IMsgHandle {
	return s.msgHandler
}

func (s *Server) Name() string {
	return s.name
}

func (s *Server) Context() context.Context {
	return s.ctx
}

func (s *Server) GetIPLimiter() ziface.IPLimiter {
	return s.safeIPLimiter
}
