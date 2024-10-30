package database

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// m *testing.M es un puntero a un objeto M del paquete testing, que proporciona métodos para ejecutar las pruebas y obtener el código de salida.
// testing.M se utiliza para controlar la ejecución de todas las pruebas dentro de un paquete. Se utiliza en la función TestMain(m *testing.M).
func TestMain(m *testing.M) {

	LoadEnv()

	initDB()

	// Ejecuta las pruebas
	code := m.Run()

	CloseDB()

	os.Exit(code)

}

// Prueba para verificar que se establece la conexión
// testing.T se utiliza dentro de las funciones de prueba individuales.
func TestInitDB(t *testing.T) {

	// Llamar a InitDB y verificar que se hace la conexión
	t.Run("Debe establecer una conexión con la base de datos", func(t *testing.T) {
		if pool == nil {
			t.Fatal("Se esperaba una conexión a la base de datos, pero fue nil")
		}

		// Verificar que la conexión está activa
		err := pool.Ping(context.Background())
		assert.NoError(t, err, "No se pudo hacer ping a la base de datos")
	})

}

// Prueba para verificar que el pool de conexiones se cierra correctamente
func TestCloseDB(t *testing.T) {
	// Llamar a CloseDB y verificar que no produce errores
	assert.NotPanics(t, func() {
		CloseDB()
	})

	// Comprobar que el pool se ha cerrado
	assert.Nil(t, pool, "Expected pool to be nil after closing")
}
