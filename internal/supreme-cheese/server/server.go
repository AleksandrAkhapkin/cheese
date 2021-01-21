package server

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/server/handlers"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"net/http"
)

func StartServer(handlers *handlers.Handlers, port string, useSSL bool, certFile, keyFile string) {
	router := NewRouter(handlers)
	if useSSL {
		if err := http.ListenAndServeTLS(port, certFile, keyFile, router); err != nil {
			logger.LogFatal(err)
		}
	}
	if err := http.ListenAndServe(port, router); err != nil {
		logger.LogFatal(err)
	}
}
