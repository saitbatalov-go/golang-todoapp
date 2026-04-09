package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var reuestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {

	if err:=json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}

	if err := reuestValidator.Struct(dest); err != nil {
		return fmt.Errorf("validate request: %w", err)
	}

	return nil
}	