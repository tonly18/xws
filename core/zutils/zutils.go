package zutils

import (
	"encoding/binary"
	"net"
	"net/http"
	"strings"
	"unsafe"

	"github.com/cespare/xxhash/v2"
	"github.com/google/uuid"
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

func GetClientIP(r *http.Request) string {
	// 1. 优先读取 X-Forwarded-For（可能有多个，用第一个）
	if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
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

// BytesToString []byte转string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes string 转[]byte
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// 生成uint64 uuid
func GenUUID64() uint64 {
	for {
		u := uuid.New()
		// 直接取高64bit
		id := binary.BigEndian.Uint64(u[:8])
		if id != 0 {
			return id
		}
		// xxhash
		if id = xxhash.Sum64(u[:]); id != 0 {
			return id
		}
	}
}
