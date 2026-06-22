package znet

import (
	"sync"
	"sync/atomic"

	"github.com/extra-time-zone/xws/core/ziface"
)

const maxConnPerIP = 100

type IPLimiter struct {
	conns sync.Map // ip => *atomic.Int32
}

func NewIPLimiter() ziface.IPLimiter {
	return &IPLimiter{}
}

func (m *IPLimiter) Allow(ip string) bool {
	v, _ := m.conns.LoadOrStore(ip, &atomic.Int32{})
	cnt := v.(*atomic.Int32)

	if cnt.Add(1) > maxConnPerIP {
		cnt.Add(-1)
		return false
	}
	return true
}

func (m *IPLimiter) Release(ip string) {
	if v, ok := m.conns.Load(ip); ok {
		if v.(*atomic.Int32).Add(-1) <= 0 {
			m.conns.Delete(ip)
		}
	}
}

func (m *IPLimiter) Range(fn func(ip, cnt any) bool) {
	m.conns.Range(fn)
}
