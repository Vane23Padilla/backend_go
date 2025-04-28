package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// ProfesoresController maneja las solicitudes relacionadas con profesores
type ProfesoresController struct {
	DB *sql.DB
}

// NewProfesoresController crea una nueva instancia del controlador de profesores
func NewProfesoresController(db *sql.DB) *ProfesoresController {
	return &ProfesoresController{DB: db}
}

// GetAllProfesores obtiene todos los profesores
func (c *ProfesoresController) GetAllProfesores(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id_, id_profesores, nombre, version FROM profesores")
	if err != nil {
		log.Printf("Error al consultar profesores: %v", err)
		http.Error(w, "Error al obtener profesores", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	profesores := []models.Profesor{}
	for rows.Next() {
		var p models.Profesor
		if err := rows.Scan(&p.ID, &p.IDProfesor, &p.Nombre, &p.Version); err != nil {
			log.Printf("Error al escanear profesor: %v", err)
			http.Error(w, "Error al procesar datos de profesores", http.StatusInternalServerError)
			return
		}
		profesores = append(profesores, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profesores)
}

// GetProfesor obtiene un profesor por su ID
func (c *ProfesoresController) GetProfesor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var p models.Profesor
	err := c.DB.QueryRow("SELECT id_, id_profesores, nombre, version FROM profesores WHERE id_profesores = ?", id).
		Scan(&p.ID, &p.IDProfesor, &p.Nombre, &p.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Profesor no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar profesor: %v", err)
		http.Error(w, "Error al obtener profesor", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
