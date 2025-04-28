package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// CiclosController maneja las solicitudes relacionadas con ciclos
type CiclosController struct {
	DB *sql.DB
}

// NewCiclosController crea una nueva instancia del controlador de ciclos
func NewCiclosController(db *sql.DB) *CiclosController {
	return &CiclosController{DB: db}
}

// GetAllCiclos obtiene todos los ciclos
func (c *CiclosController) GetAllCiclos(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id_, id_ciclos, ciclo, version FROM ciclos")
	if err != nil {
		log.Printf("Error al consultar ciclos: %v", err)
		http.Error(w, "Error al obtener ciclos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	ciclos := []models.Ciclo{}
	for rows.Next() {
		var c models.Ciclo
		if err := rows.Scan(&c.ID, &c.IDCiclo, &c.Ciclo, &c.Version); err != nil {
			log.Printf("Error al escanear ciclo: %v", err)
			http.Error(w, "Error al procesar datos de ciclos", http.StatusInternalServerError)
			return
		}
		ciclos = append(ciclos, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ciclos)
}

// GetCiclo obtiene un ciclo por su ID
func (c *CiclosController) GetCiclo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var ciclo models.Ciclo
	err := c.DB.QueryRow("SELECT id_, id_ciclos, ciclo, version FROM ciclos WHERE id_ciclos = ?", id).
		Scan(&ciclo.ID, &ciclo.IDCiclo, &ciclo.Ciclo, &ciclo.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Ciclo no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar ciclo: %v", err)
		http.Error(w, "Error al obtener ciclo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ciclo)
}
