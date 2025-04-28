package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// AsignacionesController maneja las solicitudes relacionadas con asignaciones
type AsignacionesController struct {
	DB *sql.DB
}

// NewAsignacionesController crea una nueva instancia del controlador de asignaciones
func NewAsignacionesController(db *sql.DB) *AsignacionesController {
	return &AsignacionesController{DB: db}
}

// GetAllAsignaciones obtiene todas las asignaciones
func (c *AsignacionesController) GetAllAsignaciones(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(`
		SELECT 
			pca.id_, 
			pca.id_profesores_ciclos_asignaturas, 
			pca.id_profesores, 
			pca.id_asignaturas, 
			pca.id_ciclos, 
			pca.version,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM profesores_ciclos_asignaturas pca
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
	`)
	if err != nil {
		log.Printf("Error al consultar asignaciones: %v", err)
		http.Error(w, "Error al obtener asignaciones", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	asignaciones := []models.Asignacion{}
	for rows.Next() {
		var a models.Asignacion
		if err := rows.Scan(
			&a.ID, 
			&a.IDAsignacion, 
			&a.IDProfesor, 
			&a.IDAsignatura, 
			&a.IDCiclo, 
			&a.Version,
			&a.NombreProfesor,
			&a.NombreAsignatura,
			&a.Ciclo,
		); err != nil {
			log.Printf("Error al escanear asignación: %v", err)
			http.Error(w, "Error al procesar datos de asignaciones", http.StatusInternalServerError)
			return
		}
		asignaciones = append(asignaciones, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asignaciones)
}

// GetAsignacion obtiene una asignación por su ID
func (c *AsignacionesController) GetAsignacion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var a models.Asignacion
	err := c.DB.QueryRow(`
		SELECT 
			pca.id_, 
			pca.id_profesores_ciclos_asignaturas, 
			pca.id_profesores, 
			pca.id_asignaturas, 
			pca.id_ciclos, 
			pca.version,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM profesores_ciclos_asignaturas pca
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
		WHERE pca.id_profesores_ciclos_asignaturas = ?
	`, id).Scan(
		&a.ID, 
		&a.IDAsignacion, 
		&a.IDProfesor, 
		&a.IDAsignatura, 
		&a.IDCiclo, 
		&a.Version,
		&a.NombreProfesor,
		&a.NombreAsignatura,
		&a.Ciclo,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Asignación no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar asignación: %v", err)
		http.Error(w, "Error al obtener asignación", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(a)
}

// GetAsignaturasDisponibles obtiene todas las asignaturas disponibles para matricularse
func (c *AsignacionesController) GetAsignaturasDisponibles(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(`
		SELECT 
			pca.id_profesores_ciclos_asignaturas, 
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM profesores_ciclos_asignaturas pca
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
	`)
	if err != nil {
		log.Printf("Error al consultar asignaturas disponibles: %v", err)
		http.Error(w, "Error al obtener asignaturas disponibles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	asignaturas := []map[string]interface{}{}
	for rows.Next() {
		var id, profesor, asignatura, ciclo string
		if err := rows.Scan(&id, &profesor, &asignatura, &ciclo); err != nil {
			log.Printf("Error al escanear asignatura disponible: %v", err)
			http.Error(w, "Error al procesar datos de asignaturas disponibles", http.StatusInternalServerError)
			return
		}
		asignaturas = append(asignaturas, map[string]interface{}{
			"id": id,
			"profesor": profesor,
			"asignatura": asignatura,
			"ciclo": ciclo,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(asignaturas)
}
