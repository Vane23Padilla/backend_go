package controllers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Permite CORS para desarrollo
    },
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Println("Error al establecer WebSocket:", err)
        return
    }
    defer conn.Close()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("Error al leer mensaje:", err)
            break
        }
        fmt.Printf("Mensaje recibido: %s\n", msg)

        // Puedes responder al cliente
        response := fmt.Sprintf("Servidor recibió: %s", msg)
        if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
            fmt.Println("Error al enviar mensaje:", err)
            break
        }
    }
}
