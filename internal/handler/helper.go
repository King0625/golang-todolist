package handler

type JsonResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Data    any    `json:"data"`
}
