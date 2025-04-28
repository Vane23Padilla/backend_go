package main

import (
	"log"
	"net/http"
	"os"
	"server_estudiantes/config"
	"server_estudiantes/controllers"
	"server_estudiantes/middleware"
	"server_estudiantes/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Inicializar la base de datos
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Inicializar controladores
	estudiantesController := controllers.NewEstudiantesController(db)
	asignaturasController := controllers.NewAsignaturasController(db)
	profesoresController := controllers.NewProfesoresController(db)
	ciclosController := controllers.NewCiclosController(db)
	matriculasController := controllers.NewMatriculasController(db)
	notasController := controllers.NewNotasController(db)
	asignacionesController := controllers.NewAsignacionesController(db)

	// Configurar rutas
	router := routes.SetupRoutes(
		estudiantesController,
		asignaturasController,
		profesoresController,
		ciclosController,
		matriculasController,
		notasController,
		asignacionesController,
	)

	// Agregar middleware CORS
	handler := middleware.CorsMiddleware(router)

	// Obtener puerto del servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("Servidor iniciado en http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
