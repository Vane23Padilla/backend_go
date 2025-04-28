package middleware

import (
	"net/http"
)

// CorsMiddleware agrega los encabezados CORS necesarios para permitir solicitudes desde el frontend
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir solicitudes desde cualquier origen
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// Permitir métodos HTTP específicos
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Permitir encabezados específicos
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Manejar solicitudes preflight OPTIONS
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		// Pasar al siguiente manejador
		next.ServeHTTP(w, r)
	})
}
