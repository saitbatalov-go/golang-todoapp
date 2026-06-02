package core_http_middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	core_logger "github.com/saitbatalov-go/golang-todoapp/internal/core/logger"
	core_http_response "github.com/saitbatalov-go/golang-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func CORS(allowedOrigins []string) Middleware {
    allowedOrigin := make(map[string]struct{}, len(allowedOrigins))
    
    for _, origin := range allowedOrigins {
        allowedOrigin[origin] = struct{}{}
        fmt.Printf("CORS allowed origin added: '%s'\n", origin) // ← лог
    }
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // ДЕТАЛЬНОЕ ЛОГИРОВАНИЕ
            fmt.Printf("=== CORS DEBUG ===\n")
            fmt.Printf("Request Origin: '%s'\n", origin)
            fmt.Printf("Method: %s\n", r.Method)
            fmt.Printf("Allowed origins map: %v\n", allowedOrigin)
            
            // Проверяем, разрешён ли origin
            _, isAllowed := allowedOrigin[origin]
            fmt.Printf("Is origin allowed? %v\n", isAllowed)
            
            // ВСЕГДА устанавливаем заголовки для preflight и запросов с origin
            if origin != "" {
                if isAllowed {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    fmt.Printf("Set Access-Control-Allow-Origin: %s\n", origin)
                } else {
                    fmt.Printf("WARNING: Origin '%s' NOT in allowed list\n", origin)
                    // Для отладки - временно разрешаем всё
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    fmt.Printf("DEBUG: Temporarily allowed origin: %s\n", origin)
                }
                
                w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
                w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization, X-Request-ID")
                w.Header().Set("Access-Control-Allow-Credentials", "true")
                w.Header().Set("Access-Control-Max-Age", "86400")
            }
            
            // Обработка preflight
            if r.Method == http.MethodOptions {
                fmt.Printf("Preflight request - returning 200 OK\n")
                w.WriteHeader(http.StatusOK)
                return
            }
            
            fmt.Printf("Proceeding to next handler\n")
            next.ServeHTTP(w, r)
        })
    }
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log := core_logger.FromLogger(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				">> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<< outgoing HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("duration", time.Since(before)),
			)
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctx := r.Context()
			log := core_logger.FromLogger(ctx)

			responseHanfler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if err := recover(); err != nil {

					responseHanfler.PanicResponse(err, "during handling request got unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
