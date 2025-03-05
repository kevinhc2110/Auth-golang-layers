package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	// Mantener un pool de conexiones
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	// Leer datos del .env
)

// Definir una interfaz para el pool de conexiones
type DBPool interface {
	Ping(ctx context.Context) error
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
}

// Verificar que *pgxpool.Pool implementa DBPool
var _ DBPool = (*pgxpool.Pool)(nil)

var pool DBPool

// SetDB permite inyectar una conexión (útil para pruebas)
func SetDB(mockPool DBPool) {
	pool = mockPool
}

// InitDB inicializa la conexión a la base de datos
func InitDB(connStr ...string) error {
	var databaseURL string
	if len(connStr) > 0 {
		databaseURL = connStr[0] // Usar conexión proporcionada
	} else {
		databaseURL = os.Getenv("DATABASE_URL") // Usar variable de entorno
	}

	if databaseURL == "" {
		return errors.New("la variable de entorno DATABASE_URL no está configurada")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("error al analizar la configuración de la base de datos: %w", err)
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Println("Nueva conexión establecida")
		return nil
	}

	pgPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	pool = pgPool
	log.Println("Conexión a la base de datos establecida")
	return nil
}

// GetDB devuelve el pool de conexiones
func GetDB() DBPool {
	return pool
}

// CloseDB cierra el pool de conexiones
func CloseDB() {
	if p, ok := pool.(*pgxpool.Pool); ok && p != nil {
		p.Close()
		log.Println("Pool de conexiones cerrado")
	}
	pool = nil
}
