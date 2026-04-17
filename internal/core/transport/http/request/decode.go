package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

var reuestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf(
			"decode json body:%s: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := reuestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"validate request: %v: %w",
			 err,
			 core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
