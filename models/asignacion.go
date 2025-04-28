package models

// Asignacion representa una asignación de profesor a asignatura y ciclo
type Asignacion struct {
	ID          string `json:"id_"`
	IDAsignacion string `json:"id_profesores_ciclos_asignaturas"`
	IDProfesor   string `json:"id_profesores"`
	IDAsignatura string `json:"id_asignaturas"`
	IDCiclo      string `json:"id_ciclos"`
	Version      int    `json:"version"`
	// Campos adicionales para consultas
	NombreProfesor   string `json:"nombre_profesor,omitempty"`
	NombreAsignatura string `json:"nombre_asignatura,omitempty"`
	Ciclo            string `json:"ciclo,omitempty"`
}
