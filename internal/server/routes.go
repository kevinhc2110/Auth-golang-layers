package server

import (
	"FonincoBackend/internal/server/controllers"
	"FonincoBackend/internal/server/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authControllers *controllers.AuthController) {
	// Ruta pública para iniciar sesión
	router.POST("/login", authControllers.Login)

	// Rutas protegidas con middleware de autenticación
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/change-password", authControllers.ChangePassword)

	}
}
