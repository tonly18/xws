package znet

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/tonly18/xws/core/logger"
	"github.com/tonly18/xws/core/zconf"
	"github.com/tonly18/xws/core/ziface"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WsConnectionHttpReqCtxKey http请求上下文key
type WsConnectionHttpReqCtxKey struct{}

// WsConnection 连接模块,用于处理 Websocket 连接的读写业务,
// 一个连接对应一个Connection
type WsConnection struct {
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
	property map[string]any

	// 保护当前property的锁
	propertyLock sync.Mutex

	// 当前连接的关闭状态
	isClosed bool

	// 当前连接是属于哪个Connection Manager的
	connManager ziface.IConnManager

	// 当前连接创建时Hook函数
	onConnStart func(ziface.IConnection)

	// 当前连接断开时的Hook函数
	onConnStop func(ziface.IConnection)

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

	closeCallbackMutex sync.RWMutex
}

// newServerConn 创建一个Server服务端特性的连接的方法
// Note: 名字由 NewConnection 更变
func newWebsocketConn(server ziface.IServer, conn *websocket.Conn, connID uint64, r *http.Request) ziface.IConnection {
	// 初始化Conn属性
	c := &WsConnection{
		ctx:         context.WithValue(context.Background(), WsConnectionHttpReqCtxKey{}, r.Context()),
		conn:        conn,
		connID:      connID,
		isClosed:    false,
		msgBuffChan: nil,
		property:    nil,
		name:        server.ServerName(),
		localAddr:   conn.LocalAddr().String(),
		remoteAddr:  conn.RemoteAddr().String(),
	}

	// 从server继承过来的属性
	c.packet = server.Packet()
	c.onConnStart = server.GetOnConnStart()
	c.onConnStop = server.GetOnConnStop()
	c.msgHandler = server.GetMsgHandler()

	// 将当前的Connection与Server的ConnManager绑定
	c.connManager = server.GetConnMgr()

	// 将新创建的Conn添加到连接管理中
	server.GetConnMgr().Add(c)

	return c
}

// Start 启动连接，让当前连接开始工作
func (c *WsConnection) Start() {
	ctx := c.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	c.ctx, c.cancel = context.WithCancel(ctx)

	// 按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.callOnConnStart()

	// 开启用户从客户端读取数据流程的Goroutine
	go c.StartReader()

	select {
	case <-c.ctx.Done():
		c.finalizer()
		return
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

func (c *WsConnection) GetName() string {
	return c.name
}

// Send 直接将Message数据发送数据给远程的TCP客户端
func (c *WsConnection) Send(data []byte) error {
	if c.isClosed == true {
		return errors.New("[Conn Send] connection closed when send msg")
	}

	//写回客户端: 设置写入数据流超时时间
	startTime := time.Now()
	if zconf.Config.MaxConnWriteTime > 0 {
		c.conn.SetWriteDeadline(time.Now().Add(time.Duration(zconf.Config.MaxConnWriteTime) * time.Millisecond))
	}
	if err := c.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		logger.LogErrorf("[Conn Send] writed length:%v, duration:%v, error:%v", len(data), time.Since(startTime).Milliseconds(), err)
		return err
	}

	return nil
}

// SendBuffMsg sends BuffMsg
func (c *WsConnection) SendBuffMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection closed when send buff msg")
	}

	//将data封包，并且发送
	msg, err := c.packet.Pack(NewMessage(msgID, data))
	if err != nil {
		logger.LogErrorf("[Conn SendByteMsg] Pack error msg ID = ", msgID, " Err: ", err)
		return errors.New("pack error msg")
	}

	//time out
	timer := time.NewTimer(10 * time.Millisecond)
	defer timer.Stop()

	//发送超时
	select {
	case <-c.ctx.Done():
		return errors.New("[Conn SendBuffMsg] connection closed when send buff msg")
	case <-timer.C:
		logger.LogError("[Conn SendBuffMsg] send buff msg timeout")
		return errors.New("send buff msg timeout")
	case c.msgBuffChan <- msg:
		return nil
	}

	return nil
}

func (c *WsConnection) SetProperty(key string, value any) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *WsConnection) GetProperty(key string) any {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value
	}

	return nil
}

func (c *WsConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *WsConnection) IsAlive() bool {
	if c.isClosed {
		return false
	}

	// 检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡
	return time.Now().Sub(c.lastActivityTime) < zconf.Config.HeartbeatMax
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *WsConnection) StartWriter() {
	logger.LogInfof("[Conn Writer] Writer Goroutine is running")
	defer func() {
		logger.LogInfof("[Conn Writer] %s Exit", c.GetLocalAddr())
	}()

	for {
		select {
		case data, ok := <-c.msgBuffChan:
			if ok {
				if err := c.Send(data); err != nil {
					logger.LogErrorf("[Conn Writer] Send Buff Data error, %s exit", err)
					break
				}
			} else {
				logger.LogError("[Conn Writer] msgBuffChan is Closed")
				break
			}
		case <-c.ctx.Done():
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *WsConnection) StartReader() {
	logger.LogInfo("[Conn Reader] Goroutine is running]")
	defer func() {
		c.Stop()
		logger.LogInfof("[Conn Reader] %s exit", c.GetLocalAddr())
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			messageType, data, err := c.conn.ReadMessage()
			if err != nil {
				c.cancel()
				return
			}
			if messageType == websocket.CloseMessage {
				c.cancel()
				return
			}
			if messageType == websocket.PingMessage {
				c.updateActivity()
				continue
			}
			logger.LogInfof("[Conn Reader] read buffer %s", hex.EncodeToString(data[0:len(data)]))

			//解析数据
			var msg *Message
			if err := binary.Read(bytes.NewReader(data), binary.LittleEndian, &msg); err != nil {
				logger.LogErrorf("[Conn Reader] read buffer error:%v", err)
				c.cancel()
				return
			}
			logger.LogInfof("[Conn Reader] message ID:%d, Len:%d, Data:%v", len(msg.Data), string(msg.Data))

			//Request 得到当前客户端请求的Request数据
			req := GetRequest(c, msg)

			//执行request
			if zconf.Config.WorkerPoolSize > 0 {
				c.msgHandler.SendMsgToTaskQueue(req)
			} else {
				go c.msgHandler.DoMsgHandler(req)
			}
		}
	}
}

func (c *WsConnection) finalizer() {
	// 如果用户注册了该连接的	关闭回调业务，那么在此刻应该显示调用
	c.callOnConnStop()

	c.msgLock.Lock()
	defer c.msgLock.Unlock()

	// 如果当前连接已经关闭
	if c.isClosed == true {
		return
	}

	// 关闭socket连接
	_ = c.conn.Close()

	// 将连接从连接管理器中删除
	if c.connManager != nil {
		c.connManager.Remove(c)
	}

	// 关闭该连接全部管道
	if c.msgBuffChan != nil {
		close(c.msgBuffChan)
	}

	// 设置标志位
	c.isClosed = true

	logger.LogInfof("[Conn] ConnID=%d, Conn Stop()...", c.connID)
}

func (c *WsConnection) callOnConnStart() {
	if c.onConnStart != nil {
		logger.LogInfof("[Conn] CallOnConnStart....")
		c.onConnStart(c)
	}
}

func (c *WsConnection) callOnConnStop() {
	if c.onConnStop != nil {
		logger.LogInfof("ZINX CallOnConnStop....")
		c.onConnStop(c)
	}
}

func (c *WsConnection) updateActivity() {
	c.lastActivityTime = time.Now()
}
