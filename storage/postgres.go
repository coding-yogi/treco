package storage

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/extra/pgdebug"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"log"
)

// Postgres
type Postgres struct {
	db *pg.DB
}

// GetDB
func (p Postgres) GetDB() *pg.DB {
	return p.db
}

// CreateSchema
func (p Postgres) CreateSchema(tables []interface{}) error {
	for _, t := range tables {
		err := p.db.Model(t).CreateTable(&orm.CreateTableOptions{
			Temp:          false,
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Select from postgres DB
func (p Postgres) Select(query interface{}) (interface{}, error) {
	q, ok := query.(*orm.Query)
	if !ok {
		return nil, fmt.Errorf("query should be of type orm.query")
	}

	err := q.Select()
	return nil, err
}

// Insert into postgres DB
func (p Postgres) Insert(query interface{}) error {
	q, ok := query.(*orm.Query)
	if !ok {
		return fmt.Errorf("query should be of type *orm.Query")
	}

	_, err := q.Insert()
	return err
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

	db.AddQueryHook(pgdebug.DebugHook{
		// Print all queries.
		Verbose: true,
	})

	return Postgres{db: db}, nil
}
