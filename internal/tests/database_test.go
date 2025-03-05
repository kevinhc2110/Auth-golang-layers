package tests

import (
	"FonincoBackend/internal/database"
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// mockDBPool es una implementación ficticia de DBPool para pruebas
type mockDBPool struct {
	closed bool
}

// Ping simula un ping exitoso a la base de datos
func (m *mockDBPool) Ping(ctx context.Context) error {
	if m.closed {
		return fmt.Errorf("conexión cerrada") // Simula un error cuando el pool está cerrado
	}
	return nil
}

// Close simula el cierre del pool de conexiones
func (m *mockDBPool) Close() {
	m.closed = true // Marca la conexión como cerrada
}

// Simula adquirir una conexión (devuelve nil porque no se usa en pruebas)
func (m *mockDBPool) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return nil, nil
}

// Devuelve una instancia del mock
func NewMockDB() database.DBPool {
	return &mockDBPool{}
}

// m *testing.M es un puntero a un objeto M del paquete testing, que proporciona métodos para ejecutar las pruebas y obtener el código de salida.
// testing.M se utiliza para controlar la ejecución de todas las pruebas dentro de un paquete. Se utiliza en la función TestMain(m *testing.M).
func TestMain(m *testing.M) {
	// Configura la variable de entorno para las pruebas
	os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/database_name?sslmode=disable")

	if err := database.InitDB(); err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}

	// Ejecuta las pruebas
	code := m.Run()
	os.Exit(code)

}

// Prueba para verificar que se establece la conexión
// testing.T se utiliza dentro de las funciones de prueba individuales.
func TestInitDB(t *testing.T) {

	mockDB := NewMockDB() // Usa el nuevo mock

	// Inyectar el mock en database
	database.SetDB(mockDB)

	// Llamar a InitDB y verificar que se hace la conexión
	t.Run("Debe establecer una conexión con la base de datos", func(t *testing.T) {
		db := database.GetDB()
		assert.NotNil(t, db, "Se esperaba una conexión a la base de datos, pero fue nil")

		// Simular respuesta exitosa de ping
		err := db.Ping(context.Background())
		assert.NoError(t, err, "No se pudo hacer ping a la base de datos")
	})
}

// Prueba para verificar que el pool de conexiones se cierra correctamente
func TestCloseDB(t *testing.T) {
	mockDB := NewMockDB() // Usa el nuevo mock

	// Inyectar el mock en database
	database.SetDB(mockDB)

	// Cerrar la conexión antes de la verificación
	database.CloseDB()

	// Obtener la conexión cerrada
	db := database.GetDB()

	// Verificar que `db` es nil antes de llamar Ping()
	assert.Nil(t, db, "Se esperaba que la conexión estuviera cerrada (nil), pero aún existe")

	if db != nil {
		err := db.Ping(context.Background())
		assert.Error(t, err, "Se esperaba que la conexión estuviera cerrada, pero aún responde a Ping")
	}
}
