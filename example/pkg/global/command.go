package global

// up cmd
const (
	upCmdTest = 100
)

// down cmd
const (
	downCmdTest = 100
)

// RequiredTLS 判断是否需要TLS安全验证
func RequiredTLS() bool {
	//envs := []string{ServerRunEnvTest, ServerRunEnvProd}
	//if slices.Contains(envs, *ServerRunEnv) {
	//	return true
	//}

	return true
}
