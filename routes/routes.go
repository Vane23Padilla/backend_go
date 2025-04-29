package routes

import (
	"net/http"
	"server_estudiantes/controllers"
	"server_estudiantes/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes configura todas las rutas de la API
func SetupRoutes(
	estudiantesController *controllers.EstudiantesController,
	asignaturasController *controllers.AsignaturasController,
	profesoresController *controllers.ProfesoresController,
	ciclosController *controllers.CiclosController,
	matriculasController *controllers.MatriculasController,
	notasController *controllers.NotasController,
	asignacionesController *controllers.AsignacionesController,
) http.Handler {
	router := mux.NewRouter()

	// Middleware para logging
	router.Use(middleware.LoggerMiddleware)

	// Ruta de estado
	router.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"online","server":"estudiantes"}`))
	}).Methods("GET")

	// Rutas para estudiantes
	router.HandleFunc("/estudiantes", estudiantesController.GetAllEstudiantes).Methods("GET")
	router.HandleFunc("/estudiantes/{id}", estudiantesController.GetEstudiante).Methods("GET")
	router.HandleFunc("/estudiantes", estudiantesController.CreateEstudiante).Methods("POST")
	router.HandleFunc("/estudiantes/{id}", estudiantesController.UpdateEstudiante).Methods("PUT")
	router.HandleFunc("/estudiantes/{id}", estudiantesController.DeleteEstudiante).Methods("DELETE")

	// Rutas para asignaturas
	router.HandleFunc("/api/asignaturas", asignaturasController.GetAllAsignaturas).Methods("GET")
	router.HandleFunc("/api/asignaturas/{id}", asignaturasController.GetAsignatura).Methods("GET")

	// Rutas para profesores
	router.HandleFunc("/api/profesores", profesoresController.GetAllProfesores).Methods("GET")
	router.HandleFunc("/api/profesores/{id}", profesoresController.GetProfesor).Methods("GET")

	// Rutas para ciclos
	router.HandleFunc("/api/ciclos", ciclosController.GetAllCiclos).Methods("GET")
	router.HandleFunc("/api/ciclos/{id}", ciclosController.GetCiclo).Methods("GET")

	// Rutas para asignaciones
	router.HandleFunc("/api/asignaciones", asignacionesController.GetAllAsignaciones).Methods("GET")
	router.HandleFunc("/api/asignaciones/{id}", asignacionesController.GetAsignacion).Methods("GET")

	// Rutas para matrículas
	router.HandleFunc("/api/matriculas", matriculasController.GetAllMatriculas).Methods("GET")
	router.HandleFunc("/api/matriculas", matriculasController.CreateMatricula).Methods("POST")
	router.HandleFunc("/api/matriculas/{id}", matriculasController.GetMatricula).Methods("GET")
	router.HandleFunc("/api/matriculas/{id}", matriculasController.UpdateMatricula).Methods("PUT")
	router.HandleFunc("/api/matriculas/{id}", matriculasController.DeleteMatricula).Methods("DELETE")

	// Rutas para notas
	router.HandleFunc("/api/notas", notasController.GetAllNotas).Methods("GET")
	router.HandleFunc("/api/notas/{id}", notasController.GetNota).Methods("GET")
	router.HandleFunc("/api/notas-estudiante/{id}", notasController.GetNotasByEstudiante).Methods("GET")

	// Rutas para asignaturas disponibles
	router.HandleFunc("/api/asignaturas-disponibles", asignacionesController.GetAsignaturasDisponibles).Methods("GET")
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./frontend"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return router
}
