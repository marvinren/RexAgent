package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"rexagent/pkg/conf"
	"rexagent/pkg/log"
	"time"
)

type Server struct {
	config *conf.Config
	router *mux.Router
}

func NewServer(conf *conf.Config) *Server{
	// 启动日志
	if conf.Log.Path != "" {
		log.InitLogger(conf.Log.Path)
	}
	// 配置API服务
	return &Server{
		config: conf,
	}
}



func (s *Server) Run() error{
	//设置router
	s.router = mux.NewRouter()
	//增加默认的健康见啥router
	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	s.router.HandleFunc("/command", CommandJobHandler)
	s.router.HandleFunc("/kube", KubeJobHandler)

	http_server := &http.Server{
		Addr:           s.config.Server.Listen,
		Handler:        s.router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 8192,
	}
	log.Info("Server %s running......", s.config.Server.Listen)
	err := http_server.ListenAndServe()
	if err != nil {
		log.Error("Server start error %s", err.Error())
	}
	return err
}