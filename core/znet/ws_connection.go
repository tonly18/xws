package znet

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/extra-time-zone/xws/core/logger"
	"github.com/extra-time-zone/xws/core/zconf"
	"github.com/extra-time-zone/xws/core/ziface"
	"github.com/extra-time-zone/xws/core/zutils"
	"github.com/gorilla/websocket"
)

// http请求上下文key
type WsConnHttpReqCtxKey struct{}

// WsConnection 连接模块,用于处理 Websocket 连接的读写业务,
// 一个连接对应一个Connection
type WsConnection struct {
	server  ziface.IServer
	request *http.Request
	// 当前连接的socket TCP套接字
	conn *websocket.Conn

	// 当前连接的ID 也可以称作为SessionID，ID全局唯一，服务端Connection使用
	connID uint64

	// 消息管理MsgID和对应处理方法的消息管理模块
	msgHandler ziface.IMsgHandle

	// 告知该连接已经退出/停止的channel
	ctx    context.Context
	cancel context.CancelFunc

	// 有缓冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	// 用户收发消息的Lock
	msgLock sync.Mutex

	// 连接属性
	property sync.Map

	// 当前连接的关闭状态
	isClosed bool

	// 当前连接是属于哪个Connection Manager的
	connManager ziface.IConnManager

	// 当前连接创建时Hook函数
	onConnStart func(ziface.IConnection)

	// 当前连接断开时的Hook函数
	onConnStop func(ziface.IConnection)

	// ws安全
	safeIPLimiter ziface.IPLimiter
	safeRater     ziface.IRater

	// 数据报文封包方式
	packet ziface.Packet

	// 最后一次活动时间
	lastActivityTime time.Time

	// 连接名称，默认与创建连接的Server/Client的Name一致
	name string

	// 当前连接的本地地址
	localAddr string

	// 当前连接的远程地址
	remoteAddr string

	//client ip
	clientIp string

	// 读取数据超时时长
	readDeadline time.Duration
}

// newServerConn 创建一个Server服务端特性的连接的方法
// Note: 名字由 NewConnection 更变
func newWebsocketConn(server ziface.IServer, conn *websocket.Conn, clientIp string, r *http.Request) ziface.IConnection {
	// 初始化Conn属性
	c := &WsConnection{
		server:       server,
		request:      r,
		ctx:          context.WithValue(context.Background(), WsConnHttpReqCtxKey{}, r.Context()),
		conn:         conn,
		connID:       zutils.GenUUID64(),
		isClosed:     false,
		msgBuffChan:  make(chan []byte, zconf.Config.MaxMsgChanLen),
		name:         server.Name(),
		localAddr:    conn.LocalAddr().String(),
		remoteAddr:   conn.RemoteAddr().String(),
		clientIp:     clientIp,
		readDeadline: time.Second * time.Duration(zconf.Config.HeartbeatMax*2),
	}

	// 从server继承过来的属性
	c.packet = server.Packet()
	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()
	c.msgHandler = server.GetMsgHandler()

	// 将当前的Connection与Server的ConnManager绑定
	c.connManager = server.GetConnMgr()

	//context
	c.ctx, c.cancel = context.WithCancel(context.Background())

	// WS安全
	c.safeIPLimiter = server.GetIPLimiter()
	c.safeRater = NewRater()

	// 将新创建的Conn添加到连接管理中
	server.GetConnMgr().Add(c)

	return c
}

// Start 启动连接，让当前连接开始工作
func (c *WsConnection) Start() {
	// 执行钩子方法
	c.callOnConnStart()

	// 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()
	go c.StartWriter()

	select {
	case <-c.ctx.Done():
		c.finalizer()
	case <-c.server.Context().Done():
		c.finalizer()
	}
}

// Stop 停止连接，结束当前连接状态
func (c *WsConnection) Stop() {
	c.cancel()
}

// Context 返回ctx，用于用户自定义的go程获取连接退出状态
func (c *WsConnection) Context() context.Context {
	return c.ctx
}

func (c *WsConnection) GetConnection() *websocket.Conn {
	return c.conn
}

func (c *WsConnection) GetConnID() uint64 {
	return c.connID
}
func (c *WsConnection) SetConnID(id uint64) {
	c.connID = id
}

func (c *WsConnection) GetConnMgr() ziface.IConnManager {
	return c.connManager
}

func (c *WsConnection) GetMsgHandler() ziface.IMsgHandle {
	return c.msgHandler
}

func (c *WsConnection) GetRemoteAddr() string {
	return c.remoteAddr
}

func (c *WsConnection) GetLocalAddr() string {
	return c.localAddr
}

func (c *WsConnection) GetClientIP() string {
	return c.clientIp
}

func (c *WsConnection) GetName() string {
	return c.name
}

// Send 直接将Message数据发送数据给远程的websocket客户端
func (c *WsConnection) Send(data []byte) error {
	if c.isClosed {
		return errors.New("[Conn Send] connection closed when send msg")
	}

	//deadline
	c.conn.SetWriteDeadline(time.Now().UTC().Add(time.Second * 2))

	//写回客户端: 设置写入数据流超时时间
	start := time.Now()
	if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.LogErrorf("[Conn Send] writed length:%v, duration:%v, error:%v", len(data), time.Since(start).Milliseconds(), err)
		return err
	}

	return nil
}

// Send 直接批量将Message数据发送数据给远程的websocket客户端
func (c *WsConnection) SendBatch(data []byte) error {
	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	if c.isClosed {
		return errors.New("[Conn SendBatch] connection closed when send msg")
	}

	//deadline
	c.conn.SetWriteDeadline(time.Now().UTC().Add(time.Second * 2))

	//send batch
	w, err := c.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return fmt.Errorf(`send batch next write error:%v`, w)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf(`send batch write error:%v`, w)
	}

	//loop write
	stop := false
	for i := 0; i < 50; i++ {
		select {
		case msg := <-c.msgBuffChan:
			_, err = w.Write(msg)
			if err != nil {
				return fmt.Errorf(`send batch loop write error:%v`, w)
			}
		default:
			stop = true
		}
		if stop {
			break
		}
	}
	_ = w.Close()

	return nil
}

// SendBuffMsg sends BuffMsg
func (c *WsConnection) SendBuffMsg(msgID uint32, data []byte) error {
	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	//conn 已关闭
	if c.isClosed {
		return errors.New("connection closed when send buff msg")
	}

	//deadline
	c.conn.SetWriteDeadline(time.Now().UTC().Add(time.Second * 2))

	//将data封包，并且发送
	msgRaw := msgDownGetFromPool()
	defer msgDownPutToPool(msgRaw)

	msgRaw.SetCmd(msgID)
	msgRaw.SetTime(time.Now().UTC().UnixMilli())
	msgRaw.SetData(data)
	msg, err := c.packet.Pack(msgRaw)
	if err != nil {
		logger.LogErrorf("[Conn SendBuffMsg] Pack error msg ID:%d Err:%v", msgID, err)
		return errors.New("pack error msg")
	}

	//time out
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	//发送超时
	select {
	case <-c.ctx.Done():
		return errors.New("[Conn SendBuffMsg] connection closed when send buff msg")
	case <-timer.C:
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}

	return nil
}

func (c *WsConnection) SetProperty(key, value any) {
	c.property.Store(key, value)
}

func (c *WsConnection) GetProperty(key any) any {
	if value, ok := c.property.Load(key); ok {
		return value
	}

	return nil
}

func (c *WsConnection) RemoveProperty(key any) {
	c.property.Delete(key)
}

func (c *WsConnection) IsAlive() bool {
	if c.isClosed {
		return false
	}

	// 检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡
	return time.Now().Sub(c.lastActivityTime) < c.readDeadline
}

func (c *WsConnection) GetRequest() *http.Request {
	return c.request
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *WsConnection) StartWriter() {
	logger.LogInfof("[Conn Writer] Writer Goroutine is running")
	defer c.Stop()

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if err := c.Send(data); err != nil {
					logger.LogErrorf("[Conn Writer] Send Buff Data error, %s exit", err)
					return
				}
			} else {
				logger.LogError("[Conn Writer] msgBuffChan is Closed")
				return
			}
		case <-c.ctx.Done():
			logger.LogInfof("[Conn Writer] connection writer closed, connId:%d, address:%s", c.connID, c.GetLocalAddr())
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *WsConnection) StartReader() {
	logger.LogInfo("[Conn Reader] Goroutine is running]")
	defer c.Stop()

	// 设置读取限制
	if zconf.Config.MaxPacketSize > 0 {
		c.conn.SetReadLimit(int64(zconf.Config.MaxPacketSize))
	}
	// conn存活时长
	if c.readDeadline > 0 {
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
	}

	// 收到 ping
	c.conn.SetPingHandler(func(data string) error {
		//fmt.Println("[Connection] service receive ping:", data, time.Now().Format(time.DateTime), "id:", c.connID)

		//更新conn活跃时间
		c.updateActivity()
		//设置读取超时
		_ = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
		if err := c.conn.WriteControl(websocket.PongMessage, nil, time.Now().Add(time.Second)); err != nil {
			logger.Errorf(c.ctx, "[Connection] service send pong (ID: %d), err:%v", c.connID, err)
			return err
		}
		return nil
	})

	for {
		//服务退出
		if zconf.ServiceExit.Load() {
			logger.LogInfo("[Conn Reader] service shut down...")
			time.Sleep(time.Second * 3)
			return
		}
		//如果当前连接已经关闭(1)
		if c.isClosed {
			logger.LogInfof("[Conn Reader] connection reader closed, connId:%d", c.connID)
			return
		}
		//读取message
		msgType, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			//正常关闭
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				logger.LogInfof("[Conn Reader] connection reader closed, connId:%d, error:%v", c.connID, err)
				return
			}
			//超时
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				logger.LogInfof("[Conn Reader] connection reader timeout ConnID:%d, error:%v", c.connID, err)
				return
			}

			logger.LogErrorf("[Conn Reader] message ConnID:%d error:%v", c.connID, err)
			return
		}
		//如果当前连接已经关闭(2)
		if c.isClosed {
			logger.LogInfof("[Conn Reader] connection reader closed message, connId:%d", c.connID)
			return
		}
		if msgType != websocket.BinaryMessage {
			logger.LogErrorf("[Conn Reader] message type ConnID:%d error:%v", c.connID, msgType)
			return
		}
		//package
		if len(rawMsg) < 12 {
			logger.LogErrorf("[Conn Reader] message length is zero ConnID:%d", c.connID)
			return
		}
		//rater
		if false == c.safeRater.GetLimiter().Allow() {
			logger.LogErrorf("[Conn Reader] rater limited ip:%s, ConnID:%d", c.clientIp, c.connID)
			return
		}

		//解析数据
		msg, err := NewPacket().UnPack(rawMsg)
		if err != nil {
			logger.LogErrorf("[Conn Reader] unpack ConnID:%d error:%v", c.connID, err)
			return
		}
		//logger.LogInfof("[Conn Reader] message ID:%d, data:%v", msg.GetCmd(), zutils.BytesToString(msg.GetData()))

		//更新conn活跃时间
		c.updateActivity()
		//conn存活时长
		if c.readDeadline > 0 {
			_ = c.conn.SetReadDeadline(time.Now().Add(c.readDeadline))
		}

		//Request 得到当前客户端请求的Request数据,并执行request
		req := requestGetFromPool(c, msg)
		if zconf.Config.WorkerPoolSize > 0 {
			c.msgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.msgHandler.DoMsgHandler(req)
		}
	}
}

func (c *WsConnection) finalizer() {
	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	// 当前连接是否已经关闭
	if c.isClosed {
		return
	}
	c.isClosed = true
	logger.LogInfo("[Connection] finalizer connId:", c.connID)

	// ip limiter
	c.safeIPLimiter.Release(c.clientIp)

	// websocket退出执行
	c.callOnConnStop()

	// 将连接从连接管理器中删除
	c.connManager.Remove(c)

	// 关闭该连接全部管道
	if c.msgBuffChan != nil {
		if len(c.msgBuffChan) > 0 {
			time.Sleep(time.Second * 5)
			if len(c.msgBuffChan) > 0 {
				logger.LogErrorf("[Connection] finalizer connId:%d, msgBuffChan has messages:%d", c.connID, len(c.msgBuffChan))
			}
		}
		close(c.msgBuffChan)
	}

	// 关闭socket连接
	_ = c.conn.Close()
}

func (c *WsConnection) callOnConnStart() {
	if c.onConnStart != nil {
		//logger.LogInfof("[Conn] CallOnConnStart hook connId:%d", c.connID)
		c.onConnStart(c)
	}
}

func (c *WsConnection) callOnConnStop() {
	if c.onConnStop != nil {
		//logger.LogInfof("callOnConnStop hook connId:%d", c.connID)
		c.onConnStop(c)
	}
}

func (c *WsConnection) updateActivity() {
	c.lastActivityTime = time.Now()
}
