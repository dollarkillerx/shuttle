package middlewares

import (
	"github.com/rs/cors"

	"net/http"
)

// Cors ...
func Cors() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler
}
