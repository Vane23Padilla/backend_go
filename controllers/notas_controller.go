package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"server_estudiantes/models"

	"github.com/gorilla/mux"
)

// NotasController maneja las solicitudes relacionadas con notas
type NotasController struct {
	DB *sql.DB
}

// NewNotasController crea una nueva instancia del controlador de notas
func NewNotasController(db *sql.DB) *NotasController {
	return &NotasController{DB: db}
}

// GetAllNotas obtiene todos los registros de notas
func (c *NotasController) GetAllNotas(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(`
		SELECT 
			rn.id_, 
			rn.id_registro_notas, 
			rn.id_matriculas, 
			rn.nota1, 
			rn.nota2, 
			rn.sup, 
			rn.version,
			e.nombre AS nombre_estudiante,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM registro_notas rn
		JOIN matriculas m ON rn.id_matriculas = m.id_matriculas
		JOIN estudiantes e ON m.id_estudiantes = e.id_estudiantes
		JOIN profesores_ciclos_asignaturas pca ON m.id_profesores_ciclos_asignaturas = pca.id_profesores_ciclos_asignaturas
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
	`)
	if err != nil {
		log.Printf("Error al consultar notas: %v", err)
		http.Error(w, "Error al obtener notas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notas := []models.Nota{}
	for rows.Next() {
		var n models.Nota
		if err := rows.Scan(
			&n.ID, 
			&n.IDNota, 
			&n.IDMatricula, 
			&n.Nota1, 
			&n.Nota2, 
			&n.Sup, 
			&n.Version,
			&n.NombreEstudiante,
			&n.NombreProfesor,
			&n.NombreAsignatura,
			&n.Ciclo,
		); err != nil {
			log.Printf("Error al escanear nota: %v", err)
			http.Error(w, "Error al procesar datos de notas", http.StatusInternalServerError)
			return
		}
		notas = append(notas, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notas)
}

// GetNota obtiene un registro de notas por su ID
func (c *NotasController) GetNota(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var n models.Nota
	err := c.DB.QueryRow(`
		SELECT 
			rn.id_, 
			rn.id_registro_notas, 
			rn.id_matriculas, 
			rn.nota1, 
			rn.nota2, 
			rn.sup, 
			rn.version,
			e.nombre AS nombre_estudiante,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM registro_notas rn
		JOIN matriculas m ON rn.id_matriculas = m.id_matriculas
		JOIN estudiantes e ON m.id_estudiantes = e.id_estudiantes
		JOIN profesores_ciclos_asignaturas pca ON m.id_profesores_ciclos_asignaturas = pca.id_profesores_ciclos_asignaturas
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
		WHERE rn.id_registro_notas = ?
	`, id).Scan(
		&n.ID, 
		&n.IDNota, 
		&n.IDMatricula, 
		&n.Nota1, 
		&n.Nota2, 
		&n.Sup, 
		&n.Version,
		&n.NombreEstudiante,
		&n.NombreProfesor,
		&n.NombreAsignatura,
		&n.Ciclo,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Registro de notas no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar registro de notas: %v", err)
		http.Error(w, "Error al obtener registro de notas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(n)
}

// GetNotasByEstudiante obtiene todos los registros de notas de un estudiante
func (c *NotasController) GetNotasByEstudiante(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idEstudiante := vars["id"]

	// Verificar si existe el estudiante
	var estudiante models.Estudiante
	err := c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", idEstudiante).
		Scan(&estudiante.ID, &estudiante.IDEstudiante, &estudiante.Nombre, &estudiante.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar estudiante: %v", err)
		http.Error(w, "Error al obtener notas del estudiante", http.StatusInternalServerError)
		return
	}

	rows, err := c.DB.Query(`
		SELECT 
			rn.id_, 
			rn.id_registro_notas, 
			rn.id_matriculas, 
			rn.nota1, 
			rn.nota2, 
			rn.sup, 
			rn.version,
			e.nombre AS nombre_estudiante,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM registro_notas rn
		JOIN matriculas m ON rn.id_matriculas = m.id_matriculas
		JOIN estudiantes e ON m.id_estudiantes = e.id_estudiantes
		JOIN profesores_ciclos_asignaturas pca ON m.id_profesores_ciclos_asignaturas = pca.id_profesores_ciclos_asignaturas
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
		WHERE e.id_estudiantes = ?
	`, idEstudiante)
	if err != nil {
		log.Printf("Error al consultar notas del estudiante: %v", err)
		http.Error(w, "Error al obtener notas del estudiante", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	notas := []models.Nota{}
	for rows.Next() {
		var n models.Nota
		if err := rows.Scan(
			&n.ID, 
			&n.IDNota, 
			&n.IDMatricula, 
			&n.Nota1, 
			&n.Nota2, 
			&n.Sup, 
			&n.Version,
			&n.NombreEstudiante,
			&n.NombreProfesor,
			&n.NombreAsignatura,
			&n.Ciclo,
		); err != nil {
			log.Printf("Error al escanear nota: %v", err)
			http.Error(w, "Error al procesar datos de notas", http.StatusInternalServerError)
			return
		}
		notas = append(notas, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notas)
}
