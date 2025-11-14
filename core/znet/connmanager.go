package znet

import (
	"errors"
	"github.com/tonly18/xws/core/logger"
	"github.com/tonly18/xws/core/zconf"
	"github.com/tonly18/xws/core/ziface"
	"sync"

	"github.com/spf13/cast"
)

// ConnManager 连接管理模块
type ConnManager struct {
	connections map[uint64]ziface.IConnection //map[connID]conn
	users       map[int64]uint64              //map[userId]connID
	connLock    sync.RWMutex
}

// NewConnManager 创建一个链接管理
func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		connections: make(map[uint64]ziface.IConnection, zconf.Config.MaxConn),
		users:       make(map[int64]uint64, zconf.Config.MaxConn),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connCount, userCount := connMgr.Len()

	connMgr.connLock.Lock()

	//将conn连接添加到ConnManager中
	if _, ok := connMgr.connections[conn.GetConnID()]; !ok {
		connMgr.connections[conn.GetConnID()] = conn
		connCount = len(connMgr.connections)
	}
	//如果conn(已登录),则添加到players中
	if userId := cast.ToInt64(conn.GetProperty(zconf.UserID)); userId > 0 {
		connMgr.users[userId] = conn.GetConnID()
		userCount = len(connMgr.users)
		logger.LogInfof(`[Conn Manager] Add UserID Successfully! conn number:%v, user number:%v, Address:%v`, connCount, userCount, conn.GetRemoteAddr())
	}
	connMgr.connLock.Unlock()

	logger.LogInfof(`[Conn Manager] Add Successfully! conn number:%v, user number:%v, Address:%v`, connCount, userCount, conn.GetRemoteAddr())
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID uint64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.connLock.Lock()

	//删除players
	userId := cast.ToInt64(conn.GetProperty(zconf.UserID))
	if cid, ok := connMgr.users[userId]; ok {
		if conn.GetConnID() == cid {
			delete(connMgr.users, userId)
		}
	}
	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	connMgr.connLock.Unlock()

	connCount, userCount := connMgr.Len()
	logger.LogInfof(`[Conn Manager] Remove ConnID:%v Successfully! conn number:%v, user number:%v, Address:%v`, conn.GetConnID(), connCount, userCount, conn.GetRemoteAddr())
}

// Clear 清除并停止所有连接
func (connMgr *ConnManager) Clear() {
	connMgr.connLock.Lock()

	//停止并删除全部的连接信息
	for uid, cid := range connMgr.users {
		delete(connMgr.users, uid)
		if conn, ok := connMgr.connections[cid]; ok {
			delete(connMgr.connections, cid) //删除
			conn.Stop()                      //停止
		}
	}
	//停止并删除全部的连接信息
	for cid, conn := range connMgr.connections {
		delete(connMgr.connections, cid) //删除
		conn.Stop()                      //停止
	}

	connMgr.connLock.Unlock()

	connCount, userCount := connMgr.Len()
	logger.LogInfof("[Conn Manager] Clear All Connections successfully! conn number:%v, user number:%v", connCount, userCount)
}

// GetByUid 根据userId获取链接
func (connMgr *ConnManager) GetByUid(userId int64) (ziface.IConnection, error) {
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if connID, ok := connMgr.users[userId]; ok {
		if conn, exist := connMgr.connections[connID]; exist {
			return conn, nil
		}
	}

	return nil, errors.New("connection not found")
}

// Len 获取当前连接、在线玩家数量
func (connMgr *ConnManager) Len() (int, int) {
	connMgr.connLock.RLock()
	connCount := len(connMgr.connections)
	userCount := len(connMgr.users)
	connMgr.connLock.RUnlock()

	return connCount, userCount
}
