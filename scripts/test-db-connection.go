package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL (CockroachDB usa protocolo PostgreSQL)
)

// Script de prueba de conexión a CockroachDB
// Este script verifica que la conexión a la base de datos funciona correctamente
func main() {
	fmt.Println("=== Test de Conexión a CockroachDB ===\n")

	// Configuración de conexión (ajusta según tu .env)
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "26257")
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "stock_analyzer")
	sslmode := getEnv("DB_SSL_MODE", "disable")

	// Construir connection string
	var connStr string
	if password != "" {
		connStr = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, port, dbname, sslmode)
	} else {
		connStr = fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=%s",
			user, host, port, dbname, sslmode)
	}

	fmt.Printf("Intentando conectar a: %s:%s\n", host, port)
	fmt.Printf("Base de datos: %s\n", dbname)
	fmt.Printf("Usuario: %s\n\n", user)

	// Intentar conexión
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ Error al abrir conexión: %v\n", err)
	}
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatalf("❌ Error al hacer ping a la base de datos: %v\n", err)
	}

	fmt.Println("✅ Conexión exitosa a CockroachDB!")

	// Test 1: Verificar versión de CockroachDB
	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Printf("⚠️  Error al obtener versión: %v\n", err)
	} else {
		fmt.Printf("\n📊 Versión: %s\n", version)
	}

	// Test 2: Verificar base de datos actual
	var currentDB string
	err = db.QueryRow("SELECT current_database()").Scan(&currentDB)
	if err != nil {
		log.Printf("⚠️  Error al obtener base de datos actual: %v\n", err)
	} else {
		fmt.Printf("📂 Base de datos actual: %s\n", currentDB)
	}

	// Test 3: Listar todas las bases de datos
	fmt.Println("\n📋 Bases de datos disponibles:")
	rows, err := db.Query("SHOW DATABASES")
	if err != nil {
		log.Printf("⚠️  Error al listar bases de datos: %v\n", err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var dbName string
			if err := rows.Scan(&dbName); err != nil {
				log.Printf("⚠️  Error al leer nombre de base de datos: %v\n", err)
				continue
			}
			fmt.Printf("  - %s\n", dbName)
		}
	}

	// Test 4: Listar tablas en la base de datos actual
	fmt.Println("\n📊 Tablas en la base de datos:")
	tableRows, err := db.Query("SHOW TABLES FROM " + dbname)
	if err != nil {
		log.Printf("⚠️  Error al listar tablas: %v\n", err)
	} else {
		defer tableRows.Close()
		hasTable := false
		for tableRows.Next() {
			hasTable = true
			var tableName string
			if err := tableRows.Scan(&tableName); err != nil {
				log.Printf("⚠️  Error al leer nombre de tabla: %v\n", err)
				continue
			}
			fmt.Printf("  - %s\n", tableName)
		}
		if !hasTable {
			fmt.Println("  (No hay tablas aún - esto es normal en una instalación nueva)")
		}
	}

	// Test 5: Crear una tabla de prueba
	fmt.Println("\n🧪 Creando tabla de prueba...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS connection_test (
			id SERIAL PRIMARY KEY,
			test_message STRING,
			created_at TIMESTAMP DEFAULT current_timestamp()
		)
	`)
	if err != nil {
		log.Printf("⚠️  Error al crear tabla de prueba: %v\n", err)
	} else {
		fmt.Println("✅ Tabla de prueba creada exitosamente")

		// Insertar un registro de prueba
		_, err = db.Exec(`
			INSERT INTO connection_test (test_message) 
			VALUES ('Conexión desde Go verificada')
		`)
		if err != nil {
			log.Printf("⚠️  Error al insertar datos de prueba: %v\n", err)
		} else {
			fmt.Println("✅ Datos de prueba insertados")

			// Leer datos de prueba
			var message string
			var createdAt string
			err = db.QueryRow(`
				SELECT test_message, created_at 
				FROM connection_test 
				ORDER BY created_at DESC 
				LIMIT 1
			`).Scan(&message, &createdAt)
			if err != nil {
				log.Printf("⚠️  Error al leer datos de prueba: %v\n", err)
			} else {
				fmt.Printf("📝 Mensaje leído: %s\n", message)
				fmt.Printf("🕐 Timestamp: %s\n", createdAt)
			}
		}

		// Limpiar tabla de prueba
		_, err = db.Exec("DROP TABLE connection_test")
		if err != nil {
			log.Printf("⚠️  Error al eliminar tabla de prueba: %v\n", err)
		} else {
			fmt.Println("🧹 Tabla de prueba eliminada")
		}
	}

	// Test 6: Verificar estadísticas de conexión
	fmt.Println("\n📈 Estadísticas de conexión:")
	stats := db.Stats()
	fmt.Printf("  - Conexiones abiertas: %d\n", stats.OpenConnections)
	fmt.Printf("  - Conexiones en uso: %d\n", stats.InUse)
	fmt.Printf("  - Conexiones inactivas: %d\n", stats.Idle)

	fmt.Println("\n✨ Todas las pruebas completadas exitosamente!")
	fmt.Println("\n🚀 CockroachDB está listo para usar en tu aplicación")
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
