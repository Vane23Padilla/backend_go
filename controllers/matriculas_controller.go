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

// MatriculasController maneja las solicitudes relacionadas con matrículas
type MatriculasController struct {
	DB *sql.DB
}

// NewMatriculasController crea una nueva instancia del controlador de matrículas
func NewMatriculasController(db *sql.DB) *MatriculasController {
	return &MatriculasController{DB: db}
}

// GetAllMatriculas obtiene todas las matrículas
func (c *MatriculasController) GetAllMatriculas(w http.ResponseWriter, r *http.Request) {
	rows, err := c.DB.Query(`
		SELECT 
			m.id_, 
			m.id_matriculas, 
			m.id_estudiantes, 
			m.id_profesores_ciclos_asignaturas, 
			m.version,
			e.nombre AS nombre_estudiante,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM matriculas m
		JOIN estudiantes e ON m.id_estudiantes = e.id_estudiantes
		JOIN profesores_ciclos_asignaturas pca ON m.id_profesores_ciclos_asignaturas = pca.id_profesores_ciclos_asignaturas
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
	`)
	if err != nil {
		log.Printf("Error al consultar matrículas: %v", err)
		http.Error(w, "Error al obtener matrículas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	matriculas := []models.Matricula{}
	for rows.Next() {
		var m models.Matricula
		if err := rows.Scan(
			&m.ID, 
			&m.IDMatricula, 
			&m.IDEstudiante, 
			&m.IDAsignacion, 
			&m.Version,
			&m.NombreEstudiante,
			&m.NombreProfesor,
			&m.NombreAsignatura,
			&m.Ciclo,
		); err != nil {
			log.Printf("Error al escanear matrícula: %v", err)
			http.Error(w, "Error al procesar datos de matrículas", http.StatusInternalServerError)
			return
		}
		matriculas = append(matriculas, m)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matriculas)
}

// GetMatricula obtiene una matrícula por su ID
func (c *MatriculasController) GetMatricula(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var m models.Matricula
	err := c.DB.QueryRow(`
		SELECT 
			m.id_, 
			m.id_matriculas, 
			m.id_estudiantes, 
			m.id_profesores_ciclos_asignaturas, 
			m.version,
			e.nombre AS nombre_estudiante,
			p.nombre AS nombre_profesor,
			a.nombre_asignatura,
			c.ciclo
		FROM matriculas m
		JOIN estudiantes e ON m.id_estudiantes = e.id_estudiantes
		JOIN profesores_ciclos_asignaturas pca ON m.id_profesores_ciclos_asignaturas = pca.id_profesores_ciclos_asignaturas
		JOIN profesores p ON pca.id_profesores = p.id_profesores
		JOIN asignaturas a ON pca.id_asignaturas = a.id_asignaturas
		JOIN ciclos c ON pca.id_ciclos = c.id_ciclos
		WHERE m.id_matriculas = ?
	`, id).Scan(
		&m.ID, 
		&m.IDMatricula, 
		&m.IDEstudiante, 
		&m.IDAsignacion, 
		&m.Version,
		&m.NombreEstudiante,
		&m.NombreProfesor,
		&m.NombreAsignatura,
		&m.Ciclo,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Matrícula no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al consultar matrícula: %v", err)
		http.Error(w, "Error al obtener matrícula", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

// CreateMatricula crea una nueva matrícula
func (c *MatriculasController) CreateMatricula(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IDEstudiante string `json:"id_estudiantes"`
		IDAsignacion string `json:"id_profesores_ciclos_asignaturas"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if input.IDEstudiante == "" || input.IDAsignacion == "" {
		http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
		return
	}

	// Verificar si existe el estudiante
	var estudiante models.Estudiante
	err := c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", input.IDEstudiante).
		Scan(&estudiante.ID, &estudiante.IDEstudiante, &estudiante.Nombre, &estudiante.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar estudiante: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	// Verificar si existe la asignación
	var asignacion models.Asignacion
	err = c.DB.QueryRow("SELECT id_, id_profesores_ciclos_asignaturas, id_profesores, id_asignaturas, id_ciclos, version FROM profesores_ciclos_asignaturas WHERE id_profesores_ciclos_asignaturas = ?", input.IDAsignacion).
		Scan(&asignacion.ID, &asignacion.IDAsignacion, &asignacion.IDProfesor, &asignacion.IDAsignatura, &asignacion.IDCiclo, &asignacion.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Asignación no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar asignación: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	// Verificar si ya existe la matrícula
	var count int
	err = c.DB.QueryRow("SELECT COUNT(*) FROM matriculas WHERE id_estudiantes = ? AND id_profesores_ciclos_asignaturas = ?", input.IDEstudiante, input.IDAsignacion).Scan(&count)
	if err != nil {
		log.Printf("Error al verificar matrícula existente: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "El estudiante ya está matriculado en esta asignatura", http.StatusBadRequest)
		return
	}

	// Crear matrícula
	id, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	idMatricula, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID de matrícula: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	_, err = c.DB.Exec(
		"INSERT INTO matriculas (id_, id_matriculas, id_estudiantes, id_profesores_ciclos_asignaturas, version) VALUES (?, ?, ?, ?, ?)",
		id, idMatricula, input.IDEstudiante, input.IDAsignacion, 1,
	)
	if err != nil {
		log.Printf("Error al insertar matrícula: %v", err)
		http.Error(w, "Error al crear matrícula", http.StatusInternalServerError)
		return
	}

	// Crear registro de notas
	idRegistro, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID de registro: %v", err)
		http.Error(w, "Error al crear registro de notas", http.StatusInternalServerError)
		return
	}

	idRegistroNotas, err := config.GenerateID()
	if err != nil {
		log.Printf("Error al generar ID de registro de notas: %v", err)
		http.Error(w, "Error al crear registro de notas", http.StatusInternalServerError)
		return
	}

	_, err = c.DB.Exec(
		"INSERT INTO registro_notas (id_, id_registro_notas, id_matriculas, nota1, nota2, sup, version) VALUES (?, ?, ?, ?, ?, ?, ?)",
		idRegistro, idRegistroNotas, idMatricula, 0, 0, 0, 1,
	)
	if err != nil {
		log.Printf("Error al insertar registro de notas: %v", err)
		// Eliminar la matrícula creada para mantener consistencia
		c.DB.Exec("DELETE FROM matriculas WHERE id_matriculas = ?", idMatricula)
		http.Error(w, "Error al crear registro de notas", http.StatusInternalServerError)
		return
	}

	nuevaMatricula := models.Matricula{
		ID:          id,
		IDMatricula: idMatricula,
		IDEstudiante: input.IDEstudiante,
		IDAsignacion: input.IDAsignacion,
		Version:     1,
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("CREATE", "matriculas", nuevaMatricula); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(nuevaMatricula)
}

// UpdateMatricula actualiza una matrícula existente
func (c *MatriculasController) UpdateMatricula(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var input struct {
		IDEstudiante string `json:"id_estudiantes"`
		IDAsignacion string `json:"id_profesores_ciclos_asignaturas"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	if input.IDEstudiante == "" || input.IDAsignacion == "" {
		http.Error(w, "Todos los campos son requeridos", http.StatusBadRequest)
		return
	}

	// Verificar si existe la matrícula
	var matricula models.Matricula
	err := c.DB.QueryRow("SELECT id_, id_matriculas, id_estudiantes, id_profesores_ciclos_asignaturas, version FROM matriculas WHERE id_matriculas = ?", id).
		Scan(&matricula.ID, &matricula.IDMatricula, &matricula.IDEstudiante, &matricula.IDAsignacion, &matricula.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Matrícula no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar matrícula: %v", err)
		http.Error(w, "Error al actualizar matrícula", http.StatusInternalServerError)
		return
	}

	// Verificar si existe el estudiante
	var estudiante models.Estudiante
	err = c.DB.QueryRow("SELECT id_, id_estudiantes, nombre, version FROM estudiantes WHERE id_estudiantes = ?", input.IDEstudiante).
		Scan(&estudiante.ID, &estudiante.IDEstudiante, &estudiante.Nombre, &estudiante.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar estudiante: %v", err)
		http.Error(w, "Error al actualizar matrícula", http.StatusInternalServerError)
		return
	}

	// Verificar si existe la asignación
	var asignacion models.Asignacion
	err = c.DB.QueryRow("SELECT id_, id_profesores_ciclos_asignaturas, id_profesores, id_asignaturas, id_ciclos, version FROM profesores_ciclos_asignaturas WHERE id_profesores_ciclos_asignaturas = ?", input.IDAsignacion).
		Scan(&asignacion.ID, &asignacion.IDAsignacion, &asignacion.IDProfesor, &asignacion.IDAsignatura, &asignacion.IDCiclo, &asignacion.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Asignación no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar asignación: %v", err)
		http.Error(w, "Error al actualizar matrícula", http.StatusInternalServerError)
		return
	}

	// Verificar si ya existe otra matrícula con los mismos datos
	var count int
	err = c.DB.QueryRow("SELECT COUNT(*) FROM matriculas WHERE id_estudiantes = ? AND id_profesores_ciclos_asignaturas = ? AND id_matriculas != ?", input.IDEstudiante, input.IDAsignacion, id).Scan(&count)
	if err != nil {
		log.Printf("Error al verificar matrícula existente: %v", err)
		http.Error(w, "Error al actualizar matrícula", http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, "El estudiante ya está matriculado en esta asignatura", http.StatusBadRequest)
		return
	}

	// Actualizar matrícula
	newVersion := matricula.Version + 1
	_, err = c.DB.Exec(
		"UPDATE matriculas SET id_estudiantes = ?, id_profesores_ciclos_asignaturas = ?, version = ? WHERE id_matriculas = ?",
		input.IDEstudiante, input.IDAsignacion, newVersion, id,
	)
	if err != nil {
		log.Printf("Error al actualizar matrícula: %v", err)
		http.Error(w, "Error al actualizar matrícula", http.StatusInternalServerError)
		return
	}

	matriculaActualizada := models.Matricula{
		ID:          matricula.ID,
		IDMatricula: matricula.IDMatricula,
		IDEstudiante: input.IDEstudiante,
		IDAsignacion: input.IDAsignacion,
		Version:     newVersion,
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("UPDATE", "matriculas", matriculaActualizada); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matriculaActualizada)
}

// DeleteMatricula elimina una matrícula
func (c *MatriculasController) DeleteMatricula(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Verificar si existe la matrícula
	var matricula models.Matricula
	err := c.DB.QueryRow("SELECT id_, id_matriculas, id_estudiantes, id_profesores_ciclos_asignaturas, version FROM matriculas WHERE id_matriculas = ?", id).
		Scan(&matricula.ID, &matricula.IDMatricula, &matricula.IDEstudiante, &matricula.IDAsignacion, &matricula.Version)
	if err == sql.ErrNoRows {
		http.Error(w, "Matrícula no encontrada", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error al verificar matrícula: %v", err)
		http.Error(w, "Error al eliminar matrícula", http.StatusInternalServerError)
		return
	}

	// Obtener el registro de notas asociado
	var registro models.Nota
	err = c.DB.QueryRow("SELECT id_, id_registro_notas, id_matriculas, nota1, nota2, sup, version FROM registro_notas WHERE id_matriculas = ?", id).
		Scan(&registro.ID, &registro.IDNota, &registro.IDMatricula, &registro.Nota1, &registro.Nota2, &registro.Sup, &registro.Version)
	
	// Eliminar primero el registro de notas (por la restricción de clave foránea)
	if err == nil {
		_, err = c.DB.Exec("DELETE FROM registro_notas WHERE id_matriculas = ?", id)
		if err != nil {
			log.Printf("Error al eliminar registro de notas: %v", err)
			http.Error(w, "Error al eliminar registro de notas", http.StatusInternalServerError)
			return
		}

		// Notificar al middleware sobre la eliminación del registro de notas
		if err := middleware.SendToMiddleware("DELETE", "registro_notas", registro); err != nil {
			log.Printf("Error al notificar al middleware: %v", err)
		}
	}

	// Eliminar la matrícula
	_, err = c.DB.Exec("DELETE FROM matriculas WHERE id_matriculas = ?", id)
	if err != nil {
		log.Printf("Error al eliminar matrícula: %v", err)
		http.Error(w, "Error al eliminar matrícula", http.StatusInternalServerError)
		return
	}

	// Notificar al middleware
	if err := middleware.SendToMiddleware("DELETE", "matriculas", matricula); err != nil {
		log.Printf("Error al notificar al middleware: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Matrícula eliminada correctamente"})
}
