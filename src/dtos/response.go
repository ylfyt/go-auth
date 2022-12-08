package dtos

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}
