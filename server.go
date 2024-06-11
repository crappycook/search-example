package main

import (
	"gosearcher/handlers"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func (s *HTTPServer) Run() {
	h := server.Default()

	v1 := h.Group("/api/v1")
	v1.GET("/ping", handlers.Ping)
	v1.GET("/search/book", handlers.SearchBooks)
	v1.GET("/search/book/fuzzy", handlers.FuzzySearchBooks)

	h.Spin()
}
