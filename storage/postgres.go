package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Postgres
type Postgres struct {
	db *gorm.DB
}

// GetDB
func (p Postgres) GetDB() *gorm.DB {
	return p.db
}

// Insert
func (p Postgres) Insert(model interface{}) error {
	return p.db.Create(model).Error
}

// Close db connection
func (p Postgres) Close() error {
	db, err := p.db.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func newPostgresDB(s db) (Postgres, error) {
	log.Println("connecting to Postgres")

	dsn := fmt.Sprintf("database=%v user=%v password=%v host=%v port=%v sslmode=disable", s.name, s.user, s.password, s.host, s.port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{db: db}, nil
}
