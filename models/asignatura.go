package models

// Asignatura representa una asignatura en el sistema
type Asignatura struct {
	ID            string `json:"id_"`
	IDAsignatura  string `json:"id_asignaturas"`
	Nombre        string `json:"nombre_asignatura"`
	Version       int    `json:"version"`
}
