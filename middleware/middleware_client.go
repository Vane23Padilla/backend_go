package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// NotifyOnline notifica al middleware que el servidor está en línea
func NotifyOnline() error {
	middlewareURL := os.Getenv("MIDDLEWARE_URL")
	if middlewareURL == "" {
		middlewareURL = "http://localhost:3001"
	}

	data := map[string]interface{}{
		"server":    "estudiantes",
		"status":    "online",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(middlewareURL+"/notify-online", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("Servidor notificado como en línea al middleware")
	return nil
}

// SendToMiddleware envía actualizaciones al middleware
func SendToMiddleware(operation, table string, data interface{}) error {
	middlewareURL := os.Getenv("MIDDLEWARE_URL")
	if middlewareURL == "" {
		middlewareURL = "http://localhost:3001"
	}

	payload := map[string]interface{}{
		"operation": operation,
		"table":     table,
		"data":      data,
		"source":    "estudiantes",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(middlewareURL+"/sync", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("Operación %s en tabla %s enviada al middleware", operation, table)
	return nil
}
