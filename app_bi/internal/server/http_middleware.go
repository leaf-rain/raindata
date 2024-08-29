package server

import "sync"

var respPool sync.Pool
var bufferSize = 1024

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

type middleware struct {
	*Server
}

func newMiddleware(server *Server) *middleware {
	return &middleware{
		Server: server,
	}
}
