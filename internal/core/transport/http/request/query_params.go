package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer:%v:%w ",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return &val, nil
}

func GetIntPathParams(r *http.Request, key string) (*int, error) {
	// Для пути "/users/2" достаём последний сегмент
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) == 0 {
		return nil, nil
	}

	// Берём последнюю часть
	param := parts[len(parts)-1]

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf("param='%s' not a valid integer: %w", param, err)
	}
	return &val, nil
}
