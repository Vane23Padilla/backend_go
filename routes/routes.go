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
	router.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
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
	router.HandleFunc("/asignaturas", asignaturasController.GetAllAsignaturas).Methods("GET")
	router.HandleFunc("/asignaturas/{id}", asignaturasController.GetAsignatura).Methods("GET")

	// Rutas para profesores
	router.HandleFunc("/profesores", profesoresController.GetAllProfesores).Methods("GET")
	router.HandleFunc("/profesores/{id}", profesoresController.GetProfesor).Methods("GET")

	// Rutas para ciclos
	router.HandleFunc("/ciclos", ciclosController.GetAllCiclos).Methods("GET")
	router.HandleFunc("/ciclos/{id}", ciclosController.GetCiclo).Methods("GET")

	// Rutas para asignaciones
	router.HandleFunc("/asignaciones", asignacionesController.GetAllAsignaciones).Methods("GET")
	router.HandleFunc("/asignaciones/{id}", asignacionesController.GetAsignacion).Methods("GET")

	// Rutas para matrículas
	router.HandleFunc("/matriculas", matriculasController.GetAllMatriculas).Methods("GET")
	router.HandleFunc("/matriculas", matriculasController.CreateMatricula).Methods("POST")
	router.HandleFunc("/matriculas/{id}", matriculasController.GetMatricula).Methods("GET")
	router.HandleFunc("/matriculas/{id}", matriculasController.UpdateMatricula).Methods("PUT")
	router.HandleFunc("/api/matriculas/{id}", matriculasController.DeleteMatricula).Methods("DELETE")

	// Rutas para notas
	router.HandleFunc("/notas", notasController.GetAllNotas).Methods("GET")
	router.HandleFunc("/notas/{id}", notasController.GetNota).Methods("GET")
	router.HandleFunc("/notas-estudiante/{id}", notasController.GetNotasByEstudiante).Methods("GET")

	// Rutas para asignaturas disponibles
	router.HandleFunc("/asignaturas-disponibles", asignacionesController.GetAsignaturasDisponibles).Methods("GET")
	// Servir archivos estáticos
	fs := http.FileServer(http.Dir("./frontend"))
	router.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	return router
}
