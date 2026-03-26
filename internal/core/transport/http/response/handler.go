package core_http_response

import (
	"fmt"
	"net/http"

	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
)

type HTTPResponseHandler struct {
	log *core_logger.Logger
}

func NewHTTPResponseHandler(log *core_logger.Logger) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
	}
}1

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	
	statusCode:=http.StatusInternalServerError
	err:= fmt.Errorf("unexpected error: %v",p)


}