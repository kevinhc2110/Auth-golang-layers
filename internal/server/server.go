package server

import (
	"FonincoBackend/internal/database"
	"FonincoBackend/internal/server/controllers"
	"FonincoBackend/internal/server/repositories"
	"FonincoBackend/internal/server/services"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Inicializa y arranca el servidor
func InitServer() {

	if err := godotenv.Load(); err != nil {
		log.Printf("Error al cargar el archivo .env: %v", err)
	}

	port := os.Getenv("PORT")

	// Inicializa la base de datos
	database.InitDB()
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
