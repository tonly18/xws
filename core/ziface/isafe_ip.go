package ziface

type IPLimiter interface {
	Allow(ip string) bool
	Release(ip string)
	Range(fn func(ip, conn any) bool)
}
