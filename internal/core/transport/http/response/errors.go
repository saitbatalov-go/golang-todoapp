package core_http_response

type ErrorResponse struct {
	Error   string `json:"error" example:"error message"`
	Message string `json:"message" example:"not found"`
}