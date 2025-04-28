package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// AsignaturasController maneja las solicitudes relacionadas con asignaturas
type AsignaturasController struct {
	DB *sql.DB
}

// NewAsignaturasController crea una nueva instancia del controlador de asignaturas
func NewAsignaturasController(db *sql.DB) *AsignaturasController {
	return &AsignaturasController{DB: db}
}

// GetAllAsignaturas obtiene todas las asignaturas
func (c *AsignaturasController) GetAllAsignaturas(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id_, id_asignaturas, nombre_asignatura, version FROM asignaturas")
	if err != nil {
		log.Printf("Error al consultar asignaturas: %v", err)
		http.Error(w, "Error al obtener asignaturas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	asignaturas := []models.Asignatura{}
	for rows.Next() {
		var a models.Asignatura
		if err := rows.Scan(&a.ID, &a.IDAsignatura, &a.Nombre, &a.Version); err != nil {
			log.Printf("Error al escanear asignatura: %v", err)
			http.Error(w, "Error al procesar datos de asignaturas", http.StatusInternalServerError)
			return
		}
		asignaturas = append(asignaturas, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asignaturas)
}

// GetAsignatura obtiene una asignatura por su ID
func (c *AsignaturasController) GetAsignatura(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var a models.Asignatura
	err := c.DB.QueryRow("SELECT id_, id_asignaturas, nombre_asignatura, version FROM asignaturas WHERE id_asignaturas = ?", id).
		Scan(&a.ID, &a.IDAsignatura, &a.Nombre, &a.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Asignatura no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar asignatura: %v", err)
		http.Error(w, "Error al obtener asignatura", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}
