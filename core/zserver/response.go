package zserver

// Response 响应
type Response struct {
	Code    string `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}
