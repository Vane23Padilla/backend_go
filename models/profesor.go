package models

// Profesor representa un profesor en el sistema
type Profesor struct {
	ID         string `json:"id_"`
	IDProfesor string `json:"id_profesores"`
	Nombre     string `json:"nombre"`
	Version    int    `json:"version"`
}
