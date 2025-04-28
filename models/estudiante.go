package models

// Estudiante representa un estudiante en el sistema
type Estudiante struct {
	ID          string `json:"id_"`
	IDEstudiante string `json:"id_estudiantes"`
	Nombre      string `json:"nombre"`
	Version     int    `json:"version"`
}
