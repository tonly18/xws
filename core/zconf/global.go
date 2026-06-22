package zconf

import "sync/atomic"

const (
	LocalEnv = "local"
	DevEnv   = "dev"
	TestEnv  = "test"
	ProdEnv  = "prod"

	ClientIP = "client_ip"
	TraceID  = "trace_id"
)

// 服务退出标识
var ServiceExit atomic.Bool
