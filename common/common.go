package common

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
