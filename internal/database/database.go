package database

import (
	"context"
	"log"
	"os"

	// Mantener un pool de conexiones

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	// Leer datos del .env
	"github.com/joho/godotenv"
)

var pool *pgxpool.Pool

// Cargar datos del .env
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el arcivo .env: %v", err)
	}
}

// Inicializa la conexión a la base de datos
func initDB() {

	LoadEnv()

	databaseURL := os.Getenv("DATABASE_URL")

	// Crea una configuración (config) para el pool de conexiones
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Printf("Error al analizar la configuración de la base de datos: %v", err)
	}

	// Se utiliza para ejecutar cualquier lógica adicional que desees al establecer una nueva conexión
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		log.Println("Nueva conexión establecida")
		return nil
	}

	// Conexión a la base de datos
	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	log.Println("Conexión a la base de datos establecida")
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
	}
}
