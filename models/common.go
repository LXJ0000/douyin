package models

// CommonResponse 公共的响应
type CommonResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type LoginResponse struct {
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}
