package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/config"
	"server_estudiantes/middleware"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// EstudiantesController maneja las solicitudes relacionadas con estudiantes
type EstudiantesController struct {
	DB *sql.DB
}

// NewEstudiantesController crea una nueva instancia del controlador de estudiantes
func NewEstudiantesController(db *sql.DB) *EstudiantesController {
	return &EstudiantesController{DB: db}
}

// GetAllEstudiantes obtiene todos los estudiantes
func (c *EstudiantesController) GetAllEstudiantes(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query("SELECT id_, id_estudiantes, nombre, version FROM estudiantes")
	if err != nil {
		log.Printf("Error al consultar estudiantes: %v", err)
		http.Error(w, "Error al obtener estudiantes", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	estudiantes := []models.Estudiante{}
	for rows.Next() {
		var e models.Estudiante
		if err := rows.Scan(&e.ID, &e.IDEstudiante, &e.Nombre, &e.Version); err != nil {
			log.Printf("Error al escanear estudiante: %v", err)
			http.Error(w, "Error al procesar datos de estudiantes", http.StatusInternalServerError)
			return
		}
		estudiantes = append(estudiantes, e)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estudiantes)
}

// GetEstudiante obtiene un estudiante por su ID
func (c *EstudiantesController) GetEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var e models.Estudiante
	err := c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", id).
		Scan(&e.ID, &e.IDEstudiante, &e.Nombre, &e.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar estudiante: %v", err)
		http.Error(w, "Error al obtener estudiante", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

// CreateEstudiante crea un nuevo estudiante
func (c *EstudiantesController) CreateEstudiante(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if input.Nombre == "" {
		http.Error(w, "El nombre es requerido", http.StatusBadRequest)
		return
	}

	id, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID: %v", err)
		http.Error(w, "Error al crear estudiante", http.StatusInternalServerError)
		return
	}

	idEstudiante, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID de estudiante: %v", err)
		http.Error(w, "Error al crear estudiante", http.StatusInternalServerError)
		return
	}

	_, err = c.DB.Exec(
		"INSERT INTO estudiantes (id_, id_estudiantes, nombre, version) VALUES (?, ?, ?, ?)",
		id, idEstudiante, input.Nombre, 1,
	)
	if err != nil {
		log.Printf("Error al insertar estudiante: %v", err)
		http.Error(w, "Error al crear estudiante", http.StatusInternalServerError)
		return
	}

	nuevoEstudiante := models.Estudiante{
		ID:          id,
		IDEstudiante: idEstudiante,
		Nombre:      input.Nombre,
		Version:     1,
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("CREATE", "estudiantes", nuevoEstudiante); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevoEstudiante)
}

// UpdateEstudiante actualiza un estudiante existente
func (c *EstudiantesController) UpdateEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if input.Nombre == "" {
		http.Error(w, "El nombre es requerido", http.StatusBadRequest)
		return
	}

	// Verificar si el estudiante existe
	var estudiante models.Estudiante
	err := c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", id).
		Scan(&estudiante.ID, &estudiante.IDEstudiante, &estudiante.Nombre, &estudiante.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar estudiante: %v", err)
		http.Error(w, "Error al actualizar estudiante", http.StatusInternalServerError)
		return
	}

	// Actualizar estudiante
	nuevaVersion := estudiante.Version + 1
	_, err = c.DB.Exec(
		"UPDATE estudiantes SET nombre = ?, version = ? WHERE id_estudiantes = ?",
		input.Nombre, nuevaVersion, id,
	)
	if err != nil {
		log.Printf("Error al actualizar estudiante: %v", err)
		http.Error(w, "Error al actualizar estudiante", http.StatusInternalServerError)
		return
	}

	estudianteActualizado := models.Estudiante{
		ID:          estudiante.ID,
		IDEstudiante: estudiante.IDEstudiante,
		Nombre:      input.Nombre,
		Version:     nuevaVersion,
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("UPDATE", "estudiantes", estudianteActualizado); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(estudianteActualizado)
}

// DeleteEstudiante elimina un estudiante
func (c *EstudiantesController) DeleteEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Verificar si el estudiante existe
	var estudiante models.Estudiante
	err := c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", id).
		Scan(&estudiante.ID, &estudiante.IDEstudiante, &estudiante.Nombre, &estudiante.Version)

	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar estudiante: %v", err)
		http.Error(w, "Error al eliminar estudiante", http.StatusInternalServerError)
		return
	}

	// Verificar si el estudiante tiene matrículas
	var count int
	err = c.DB.QueryRow("SELECT COUNT(*) FROM matriculas WHERE id_estudiantes = ?", id).Scan(&count)
	if err != nil {
		log.Printf("Error al verificar matrículas: %v", err)
		http.Error(w, "Error al eliminar estudiante", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "No se puede eliminar el estudiante porque tiene matrículas", http.StatusBadRequest)
		return
	}

	// Eliminar estudiante
	_, err = c.DB.Exec("DELETE FROM estudiantes WHERE id_estudiantes = ?", id)
	if err != nil {
		log.Printf("Error al eliminar estudiante: %v", err)
		http.Error(w, "Error al eliminar estudiante", http.StatusInternalServerError)
		return
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("DELETE", "estudiantes", estudiante); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Estudiante eliminado correctamente"})
}
