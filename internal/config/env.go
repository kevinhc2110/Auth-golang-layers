package config

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnv carga las variables de entorno desde el archivo .env
func LoadEnv() error {
	err := godotenv.Load() // Carga el archivo .env
	if err != nil {
		log.Println("No se pudo cargar el archivo .env")
		return err // Retorna el error si la carga falla
	}
	return nil // Retorna nil si no hay errores
}
