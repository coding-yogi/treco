package storage

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"log"
)

// Postgres
type Postgres struct {
	db *pg.DB
}

// Schema creation
func (p Postgres) Schema(tables []interface{}) error {
	for _, t := range tables {
		err := p.db.Model(t).CreateTable(&orm.CreateTableOptions{
			Temp: false,
			IfNotExists: true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert into postgres DB
func (p Postgres) Insert(query interface{}) error {
	if _, err := p.db.Model(query).Insert(); err != nil {
		return err
	}
	return nil
}

// Close db connection
func (p Postgres) Close() error {
	return p.db.Close()
}

func newPostgresDB(s db) (Postgres, error) {
	log.Println("connecting to Postgres")

	db := pg.Connect(&pg.Options{
		Addr:     s.host + ":" + s.port,
		User:     s.user,
		Password: s.password,
		Database: s.name,
	})

	log.Println("checking if db is up and running")
	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		return Postgres{}, err
	}

	return Postgres{db: db}, nil
}
