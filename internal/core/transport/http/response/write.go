package core_http_response

import "net/http"

type ResponseWrite struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter)*ResponseWrite {
	return &ResponseWrite{}
}