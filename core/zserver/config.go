package zserver

import (
	"net/http"
)

type HttpServerConfig struct {
	IP      string
	Port    int
	Handler http.Handler
}
