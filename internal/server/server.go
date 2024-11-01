package server

import (
	"FonincoBackend/internal/config"
	"FonincoBackend/internal/database"
	"FonincoBackend/internal/server/controllers"
	"FonincoBackend/internal/server/repositories"
	"FonincoBackend/internal/server/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// Inicializa y arranca el servidor
func InitServer() {

	// Cargar configuraciones del entorno
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error al cargar las variables de entorno: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Puerto por defecto
	}

	// Inicializa la base de datos
	if err := database.InitDB(); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer database.CloseDB() // Asegúrate de cerrar la conexión al final

	router := gin.Default()

	// Crear el repositorio y servicios
	userRepo := repositories.NewAuthRepository(database.GetDB())
	authService := &services.AuthService{UserRepo: userRepo}
	authController := controllers.NewAuthController(authService, userRepo)

	// Registrar rutas
	RegisterRoutes(router, authController)

	// Inicia el servidor
	log.Printf("Iniciando el servidor en %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
