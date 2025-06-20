package models

type LoginRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Id   string `json:"id"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
