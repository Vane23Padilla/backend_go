package models

// Matricula representa una matrícula de un estudiante en una asignación
type Matricula struct {
	ID          string `json:"id_"`
	IDMatricula string `json:"id_matriculas"`
	IDEstudiante string `json:"id_estudiantes"`
	IDAsignacion string `json:"id_profesores_ciclos_asignaturas"`
	Version      int    `json:"version"`
	// Campos adicionales para consultas
	NombreEstudiante string `json:"nombre_estudiante,omitempty"`
	NombreProfesor   string `json:"nombre_profesor,omitempty"`
	NombreAsignatura string `json:"nombre_asignatura,omitempty"`
	Ciclo            string `json:"ciclo,omitempty"`
}
