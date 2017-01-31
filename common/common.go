package common

import (
	"net/http"
	"time"
)

type ApiResponse struct {
	Status string      `json:"status"`
	Resp   interface{} `json:"data"`
}

func NewApiResponse(r interface{}) *ApiResponse {
	if t, ok := r.(error); ok {
		return &ApiResponse{Status: "error", Resp: t.Error()}
	}
	return &ApiResponse{Status: "ok", Resp: r}
}

type History struct {
	Duration time.Duration
	MinBytes int
}

type Context struct {
	History map[string]History
	Client  map[string]*http.Client
}

func NewContext() Context {
	return Context{make(map[string]History), make(map[string]*http.Client)}
}
