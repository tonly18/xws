package znet

import (
	"errors"
	"sync"

	"github.com/extra-time-zone/xws/core/logger"
	"github.com/extra-time-zone/xws/core/zconf"
	"github.com/extra-time-zone/xws/core/ziface"
)

// ConnManager 连接管理模块
type ConnManager struct {
	lock  sync.RWMutex
	conns map[uint64]ziface.IConnection //map[connId]conn
}

// NewConnManager 创建一个链接管理
func NewConnManager() ziface.IConnManager {
	return &ConnManager{
		conns: make(map[uint64]ziface.IConnection, zconf.Config.MaxConn),
	}
}

// Add 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	connMgr.lock.Lock()
	defer connMgr.lock.Unlock()

	if _, ok := connMgr.conns[conn.GetConnID()]; ok {
		logger.LogErrorf(`[Conn Manager] conn already exists, Address:%v`, conn.GetRemoteAddr())
		return
	}
	connMgr.conns[conn.GetConnID()] = conn

	logger.LogInfof(`[Conn Manager] Add Successfully! Address:%v`, conn.GetRemoteAddr())
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) Get(connId uint64) (ziface.IConnection, error) {
	connMgr.lock.RLock()
	defer connMgr.lock.RUnlock()

	if conn, ok := connMgr.conns[connId]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) GetConns(connIds []uint64) []ziface.IConnection {
	connMgr.lock.RLock()
	defer connMgr.lock.RUnlock()

	conns := make([]ziface.IConnection, 0, len(connIds))
	for _, connId := range connIds {
		if conn, ok := connMgr.conns[connId]; ok {
			conns = append(conns, conn)
		}
	}

	return conns
}

// Get 利用ConnID获取链接
func (connMgr *ConnManager) GetAllConns() []ziface.IConnection {
	connMgr.lock.RLock()
	defer connMgr.lock.RUnlock()

	conns := make([]ziface.IConnection, 0, len(connMgr.conns))
	for _, c := range connMgr.conns {
		conns = append(conns, c)
	}

	return conns
}

// RangeConn 安全遍历Conn
func (connMgr *ConnManager) RangeConn(fn func(ziface.IConnection)) {
	connMgr.lock.RLock()
	for _, conn := range connMgr.conns {
		fn(conn)
	}
	connMgr.lock.RUnlock()
}

// Remove 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	connMgr.lock.Lock()
	delete(connMgr.conns, conn.GetConnID())
	connMgr.lock.Unlock()

	logger.LogInfof(`[Conn Manager] Remove Conn Successfully! ConnID:%v, ClientIP:%v, Address:%v`, conn.GetConnID(), conn.GetClientIP(), conn.GetRemoteAddr())
}

// Clear 清除并停止所有连接
func (connMgr *ConnManager) Clear() {
	connMgr.lock.RLock()
	//停止并删除全部的连接信息
	for _, conn := range connMgr.conns {
		conn.Stop()
	}
	connMgr.lock.RUnlock()

	logger.LogInfo("[Conn Manager] Clear All Connections successfully!")
}

// ConnCount 获取当前连接数量
func (connMgr *ConnManager) Count() int {
	connMgr.lock.RLock()
	defer connMgr.lock.RUnlock()

	return len(connMgr.conns)
}
