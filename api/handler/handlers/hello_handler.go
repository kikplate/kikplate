package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kickplate/api/lib"
)

type HelloHandler struct {
	logger lib.Logger
	env    lib.Env
}

func NewHelloHandler(
	logger lib.Logger,
	env lib.Env,
) HelloHandler {
	return HelloHandler{
		logger: logger,
		env:    env,
	}
}

func (h HelloHandler) Hello(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request for /hello")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "Hello, World!",
		"env":      h.env.Environment,
		"logLevel": h.env.LogLevel,
	})
}
