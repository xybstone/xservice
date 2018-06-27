package service

import (
	"net"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/xybstone/fasthttp-routing"
)

var (
	// ServiceName 服务名称
	ServiceName = "XService"
	// ServiceVersion 版本
	ServiceVersion = "1.0.0.1"
	// Concurrency Concurrency
	Concurrency = 100000
	// DisableKeepalive DisableKeepalive
	DisableKeepalive = true
	//Timeout Timeout
	Timeout = 60 * time.Second
	// MaxConnsPerIP MaxConnsPerIP
	MaxConnsPerIP = 100000
	// MaxRequestsPerConn MaxRequestsPerConn
	MaxRequestsPerConn = 100000
	// MaxKeepaliveDuration MaxKeepaliveDuration
	MaxKeepaliveDuration = 120 * time.Second
	// MaxRequestBodySize MaxRequestBodySize
	MaxRequestBodySize = 512 * 1024 * 1024
	// ReadBufferSize ReadBufferSize
	ReadBufferSize = 16 * 1024
	// WriteBufferSize WriteBufferSize
	WriteBufferSize = 16 * 1024
)

// Logger 日志插件
var XLogger fasthttp.Logger

func GetRoute() *routing.Router {
	router := routing.New()
	RegsitRouter(router, new(BaseServer))
	return router
}

func GetServer() *fasthttp.Server {
	r := GetRoute()
	if r != nil {
		return &fasthttp.Server{
			Handler:              r.HandleRequest,
			Name:                 ServiceName,
			Concurrency:          Concurrency,
			DisableKeepalive:     DisableKeepalive,
			ReadTimeout:          Timeout,
			WriteTimeout:         Timeout,
			MaxConnsPerIP:        MaxConnsPerIP,
			MaxRequestsPerConn:   MaxRequestsPerConn,
			MaxKeepaliveDuration: MaxKeepaliveDuration,
			MaxRequestBodySize:   MaxRequestBodySize,
			ReadBufferSize:       ReadBufferSize,
			WriteBufferSize:      WriteBufferSize,
			Logger:               XLogger,
			LogAllErrors:         true,
		}
	}
	return nil
}

// Run 启动服务
func Run(n net.Listener) {
	s := GetServer()
	if s != nil {
		s.Serve(n)
	}

}
