package repositories

import (
	"FonincoBackend/internal/database"
	"FonincoBackend/internal/server/models"
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Define los metodos para manejar la autenticacion de los usuarios
type AuthRepository interface {
	FindUserByID(userID string) (*models.LoginRequest, error)
	ChangePassword(userID, newPassword string) error
}

// Implementacion concreta de AuthRepository
type authRepository struct {
	pool *pgxpool.Pool
}

// NewAuthRepository crea una nueva instancia de AuthRepository
func NewAuthRepository(pool database.DBPool) AuthRepository {
	pgxPool, ok := pool.(*pgxpool.Pool)
	if !ok {
		log.Fatal("Error: El pool no es un *pgxpool.Pool")
	}
	return &authRepository{pool: pgxPool}
}

// FindUserByID busca un usuario por su ID
func (r *authRepository) FindUserByID(userID string) (*models.LoginRequest, error) {
	var user models.LoginRequest
	query := `SELECT user_id, password FROM users WHERE user_id = $1`
	err := r.pool.QueryRow(context.Background(), query, userID).Scan(&user.UserID, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Usuario no encontrado en la base de datos")
			return nil, nil
		}
		log.Println("Error en la consulta:", err)
		return nil, err
	}
	return &user, nil
}

// ChangePassword cambia la contraseña del usuario
func (r *authRepository) ChangePassword(userID, newPassword string) error {
	query := "UPDATE users SET password = $1 WHERE user_id = $2"
	_, err := r.pool.Exec(context.Background(), query, newPassword, userID)
	return err
}
