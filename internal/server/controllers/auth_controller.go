package controllers

import (
	"FonincoBackend/internal/server/models"
	"FonincoBackend/internal/server/repositories"
	"FonincoBackend/internal/server/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Estructura que contiene los servicios y repositorios necesarios
type AuthController struct {
	AuthService *services.AuthService
	UserRepo    repositories.AuthRepository
}

// Constructor que inicializa un nuevo AuthController
func NewAuthController(authService *services.AuthService, userRepo repositories.AuthRepository) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserRepo:    userRepo,
	}
}

// Login
func (ac *AuthController) Login(c *gin.Context) {
	var credentials models.LoginRequest

	// Intenta enlazar el cuerpo JSON de la solicitud al modelo de credenciales de inicio de sesión
	if err := c.ShouldBindJSON(&credentials); err != nil {

		// Si el JSON no es válido o está incompleto, responde con un error 400 (Bad Request)
		log.Printf("Usuario no encontrado: %v", credentials.UserID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrada invalida"})
		return
	}

	// Intenta autenticar al usuario con las credenciales proporcionadas
	token, err := ac.AuthService.LoginUser(credentials.UserID, credentials.Password)
	if err != nil {
		log.Printf("Error de autenticación para usuario %s: %v", credentials.UserID, err)

		// Si la autenticación falla, responde con un error 401 (Unauthorized)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Permite al usuario cambiar su contraseña despues del primer login
func (ac *AuthController) ChangePassword(c *gin.Context) {
	var ChangePasswordRequest models.ChangePasswordRequest

	// Intenta enlazar el cuerpo JSON de la solicitud al modelo de cambio de contraseña
	if err := c.ShouldBindJSON(&ChangePasswordRequest); err != nil {

		// Si el JSON no es válido o está incompleto, responde con un error 400 (Bad Request)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Entrada invalida"})
		return
	}

	// Intenta cambiar la contraseña del usuario
	if err := ac.AuthService.ChangePassword(ChangePasswordRequest.UserID, ChangePasswordRequest.OldPassword, ChangePasswordRequest.NewPassword); err != nil {

		// Si la contraseña no se puede cambiar, responde con un error 400 (Bad Request)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cambio de contraseña exitoso"})
}
