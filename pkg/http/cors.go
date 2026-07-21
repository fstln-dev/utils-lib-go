package http

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

func Cors(origins, methods, headers []string) http.ServerOption {
	return http.Filter(
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedOrigins(origins),
			handlers.AllowedMethods(methods),
			handlers.AllowedHeaders(headers),
		),
	)
}
