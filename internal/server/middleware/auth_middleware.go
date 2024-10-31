package middleware

import (
	"errors"
	"fmt"

	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Clave secreta para firmar el token
var secretKey = os.Getenv("JWT_SECRET_KEY")

// Middleware para validar el token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Clave secreta leída del entorno:", secretKey)
		// Obtener el token de las cabeceras de la solicitud
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta el encabezado de autorización"})
			c.Abort()
			return
		}

		// Eliminar "Bearer " del token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta token"})
			c.Abort()
			return
		}

		// Verificar el token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar que el método de firma sea HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de firma inesperado")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			fmt.Println("Error al parsear el token:", err)
			fmt.Println("Clave secreta utilizada para verificar:", secretKey)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("Token no válido")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido"})
			c.Abort()
			return
		}

		// Extraer el ID de usuario (o cualquier otro campo) del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Asegurarse de que user_id está presente en las claims
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id no encontrado en claims"})
			c.Abort()
			return
		}

		c.Set("userID", userID) // Añadir userID al contexto para acceso en controladores

		// Continúa al siguiente controlador
		c.Next()
	}
}
