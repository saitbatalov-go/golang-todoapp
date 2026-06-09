package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/saitbatalov-go/golang-todoapp/internal/core/errors"
)

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := uuid.Parse(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid uuid: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
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

func GetTimeQueryParams(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	val, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid time:%v:%w ",
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

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"

	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}