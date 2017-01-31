package common

import (
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
}

func NewContext() Context {
	return Context{make(map[string]History)}
}
