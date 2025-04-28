package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware registra información sobre cada solicitud HTTP
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Llamar al siguiente manejador
		next.ServeHTTP(w, r)
		
		// Registrar información de la solicitud
		log.Printf(
			"%s %s %s",
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
