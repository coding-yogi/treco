package storage

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres DB
type Postgres struct {
	db *gorm.DB
}

// Setup creates the required entities in DB
func (p Postgres) Setup(entities ...interface{}) error {
	return p.db.AutoMigrate(entities...)
}

// GetDB returns DB instance
func (p Postgres) GetDB() *gorm.DB {
	return p.db
}

// Insert model into DB
func (p Postgres) Insert(model interface{}) error {
	return p.db.Create(model).Error
}

// Close DB connection
func (p Postgres) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

var connectToPostgresDB = func(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func newPostgresDB(s db) (Postgres, error) {
	log.Println("connecting to Postgres")

	dsn := fmt.Sprintf("database=%v user=%v password=%v host=%v port=%v", s.Name, s.User, s.Password, s.Host, s.Port)
	db, err := connectToPostgresDB(dsn)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{db: db}, nil
}
