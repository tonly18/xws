package zutils

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"net"
	"net/http"
	"strings"
	"unsafe"

	"github.com/spf13/cast"
)

type eface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

// IsNil 值判空
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	ep := (*eface)(unsafe.Pointer(&v))
	if ep == nil {
		return true
	}

	return ep.typ == nil || uintptr(ep.ptr) == 0x0
}

// GenTraceID 生成链路追踪ID
func GenTraceID() string {
	randomNum, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))

	return cast.ToString(randomNum)
}

func GetClientIP(r *http.Request) string {
	// 1. 优先读取 X-Forwarded-For（可能有多个，用第一个）
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		// 第一段就是客户端 IP
		ip := strings.TrimSpace(parts[0])
		if net.ParseIP(ip) != nil {
			return ip
		}
	}

	// 2. 再读 X-Real-IP
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		if net.ParseIP(xr) != nil {
			return xr
		}
	}

	// 3. 最后兜底 RemoteAddr（一般适用于无代理或反代未设置）
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && net.ParseIP(ip) != nil {
		return ip
	}

	return ""
}
