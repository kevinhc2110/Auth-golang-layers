package services

import (
	"FonincoBackend/internal/server/repositories"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo repositories.AuthRepository
}

// Verifica la cedula y contraseña del usuario
func (as *AuthService) LoginUser(userID string, password string) (string, error) {
	// Busca el usuario por cedula en el repositorio
	user, err := as.UserRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return "", errors.New("usuario no encontrado")
	}

	// Verifica que la contraseña coincida usando bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("credenciales invalidas")
	}

	// Generar el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Expira en 24 horas
	})

	// Firmar el token con la clave secreta
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("clave secreta no configurada")
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	fmt.Println("Token generado:", tokenString)
	return tokenString, nil
}

// Cambia la contraseña de un usuario
func (as *AuthService) ChangePassword(userID string, oldPassword, newPassword string) error {
	user, err := as.UserRepo.FindUserByID(userID)
	if err != nil || user == nil {
		return errors.New("usuario no encontrado")
	}

	// Verifica que la contraseña antigua coincida
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("contraseña antigua incorrecta")
	}

	// Genera el hash para la nueva contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	/// Actualiza la contraseña usando el repositorio
	return as.UserRepo.ChangePassword(userID, string(hashedPassword))
}
