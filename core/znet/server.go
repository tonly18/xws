package znet

import (
	"context"
	"fmt"
	"github.com/tonly18/xws/core/logger"
	"github.com/tonly18/xws/core/zconf"
	"github.com/tonly18/xws/core/ziface"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Server struct {
	Env string
	//服务器ID
	ID uint32
	//服务器的名称
	Name string
	//服务绑定的IP地址
	IP string
	//服务绑定的端口
	Port int
	//path
	Path string
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

	ctx    context.Context
	cancel context.CancelFunc
	packet ziface.Packet
	connId uint64
}

// NewServer 创建一个服务器句柄
func NewServer() ziface.IServer {
	s := &Server{
		Env:        zconf.Config.ENV,
		ID:         zconf.Config.ServerID,
		Name:       zconf.Config.Name,
		IP:         zconf.Config.Host,
		Port:       zconf.Config.Port,
		Path:       zconf.Config.Path,
		msgHandler: NewMsgHandle(),
		connMgr:    NewConnManager(),
		packet:     NewPacket(),
		ctx:        nil,
		cancel:     nil,
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
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
	logger.LogInfof("[WebSocket Server] Successful Server Name:%s, Listen At:%s:%v, Path:%v", s.Name, s.IP, s.Port, s.Path)

	http.HandleFunc(s.Path, func(w http.ResponseWriter, r *http.Request) {
		//1. 设置服务器最大连接控制,如果超过最大连接，则等待
		if connCount, _ := s.GetConnMgr().Len(); connCount >= zconf.Config.MaxConn {
			logger.LogInfof("[WebSocket Server] Exceeded the maxConnNum:%d, Wait:%d", zconf.Config.MaxConn, AcceptDelay.duration)
			AcceptDelay.Delay()
			return
		}

		//2. 判断 header 里面是有子协议
		if len(r.Header.Get("Sec-Websocket-Protocol")) > 0 {
			s.upgrader.Subprotocols = websocket.Subprotocols(r)
		}

		//3. 升级成 websocket 连接
		conn, err := s.upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.LogErrorf("[WebSocket Server] New Websocket error:%v", err)
			w.WriteHeader(500)
			AcceptDelay.Delay()
			return
		}
		AcceptDelay.Reset()

		//4. 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
		cid := atomic.AddUint64(&s.connId, 1)
		wsConn := newWebsocketConn(s, conn, cid, r)

		go s.StartConn(wsConn)
	})

	if zconf.Config.CertFile != "" && zconf.Config.PrivateKeyFile != "" {
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", s.IP, s.Port), zconf.Config.CertFile, zconf.Config.PrivateKeyFile, nil); err != nil {
			panic(err)
		}
	} else {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.IP, s.Port), nil); err != nil {
			panic(err)
		}
	}
}

// Stop 停止服务
func (s *Server) Stop() {
	logger.LogInfo("[WebSocket Server] Server Name:", s.Name)

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
	return s.ID
}

// GetMsgHandler 获取Server绑定的消息处理模块
func (s *Server) GetMsgHandler() ziface.IMsgHandle {
	return s.msgHandler
}

func (s *Server) ServerName() string {
	return s.Name
}
