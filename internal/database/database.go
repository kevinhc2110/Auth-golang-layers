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

var pool *pgxpool.Pool

// Inicializa la conexión a la base de datos
func InitDB() error {

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return errors.New("la variable de entorno DATABASE_URL no está configurada")
	}

	// Crea una configuración (config) para el pool de conexiones
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("error al analizar la configuración de la base de datos: %w", err)
	}

	// Se utiliza para ejecutar cualquier lógica adicional que desees al establecer una nueva conexión
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Println("Nueva conexión establecida")
		return nil
	}

	// Conexión a la base de datos
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("error al conectar a la base de datos: %w", err)
	}
	log.Println("Conexión a la base de datos establecida")
	return nil
}

// Devuelve el pool de conexiones
func GetDB() *pgxpool.Pool {
	return pool
}

// Cerrar conexiones
func CloseDB() {
	if pool != nil {
		pool.Close()
		log.Println("Pool de conexiones cerrado")
		pool = nil
	}
}
