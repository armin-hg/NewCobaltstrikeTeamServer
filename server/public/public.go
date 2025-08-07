package public

type ApiResponse struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Count int64       `json:"count,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}
