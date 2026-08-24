// Package server 提供 gauss-plume 的薄 HTTP 层：
//
//	POST /api/conc    单受体浓度核算
//	POST /api/axis    下风向轴线核算
//	GET  /api/version 版本信息
//	GET  /example/*   离线算例文件
//	GET  /*           静态前端（web/）
//
// 所有非法输入都通过 JSON 错误体 {"error": "..."} 返回，绝不静默给数值。
package server

import (
	"io/fs"
	"net/http"
)

// Server 持有路由与静态资源。
type Server struct {
	webFS     fs.FS
	exampleFS fs.FS
	version   string
}

// New 构造 Server。webFS 提供前端页面，exampleFS 提供离线算例。
func New(webFS, exampleFS fs.FS, version string) *Server {
	return &Server{webFS: webFS, exampleFS: exampleFS, version: version}
}

// Handler 组装全部路由。具体路由在前，通配 "/" 兜底静态文件。
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/conc", s.handleConc)
	mux.HandleFunc("/api/axis", s.handleAxis)
	mux.HandleFunc("/api/version", s.handleVersion)
	mux.Handle("/example/", http.StripPrefix("/example/", http.FileServer(http.FS(s.exampleFS))))
	mux.Handle("/", http.FileServer(http.FS(s.webFS)))
	return mux
}

// Version 返回构建时注入的版本号。
func (s *Server) Version() string { return s.version }
