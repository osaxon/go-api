package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-api/internal/dal"
	"go-api/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

var Q *dal.Query

// DAO represents a Service that interacts with a database.
type DAO interface {
	// Health returns a map of health status information.
	// The keys and values in the map are Service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	// GetDB gets the generic database interface *sql.DB from the current *gorm.DB
	getDB() *sql.DB
	GetORM() *gorm.DB
}

type Service struct {
	*gorm.DB
}

var (
	dbUrl      = os.Getenv("DB_URL")
	dbInstance *Service
)

func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.Task{})
}

func New() *Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	//db, err := sql.Open("sqlite3", dbUrl)
	db, err := gorm.Open(sqlite.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//dal.Use(db)

	dbInstance = &Service{
		db,
	}
	migrate(dbInstance.DB)
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *Service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	//db := s.GetDB()
	db := s.getDB()
	err := db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *Service) Close() error {
	log.Printf("Disconnected from database: %s", dbUrl)
	db := s.getDB()
	return db.Close()
}

func (s *Service) getDB() *sql.DB {
	sqlDb, err := s.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get DB from database: %v", err)
		return nil
	}
	return sqlDb
}

func (s *Service) GetORM() *gorm.DB {
	return s.DB
}
