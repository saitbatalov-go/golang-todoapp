package core_http_response

import "net/http"

const (
	StatusCodeUninitialized = -1
)

type ResponseWrite struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter)*ResponseWrite {
	return &ResponseWrite{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWrite) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
} 

func (rw *ResponseWrite) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}

	return rw.statusCode
}