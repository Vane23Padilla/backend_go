package models

// Nota representa un registro de notas de un estudiante
type Nota struct {
	ID        string  `json:"id_"`
	IDNota    string  `json:"id_registro_notas"`
	IDMatricula string `json:"id_matriculas"`
	Nota1     float64 `json:"nota1"`
	Nota2     float64 `json:"nota2"`
	Sup       int     `json:"sup"`
	Version   int     `json:"version"`
	// Campos adicionales para consultas
	NombreEstudiante string `json:"nombre_estudiante,omitempty"`
	NombreProfesor   string `json:"nombre_profesor,omitempty"`
	NombreAsignatura string `json:"nombre_asignatura,omitempty"`
	Ciclo            string `json:"ciclo,omitempty"`
}
