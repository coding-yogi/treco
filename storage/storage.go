package storage

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DBHost     = "DB_HOST"
	DBPort     = "DB_PORT"
	DBName     = "DB_NAME"
	DBUser     = "DB_USER"
	DBPassword = "DB_PASSWORD"
	DBType     = "DB_TYPE"
)

var dbHandler DBHandler

// DBHandler interface
type DBHandler interface {
	Insert(model interface{}) error
	Close() error
}

type db struct {
	dbType   string
	host     string
	port     string
	name     string
	user     string
	password string
}

func Handler() *DBHandler {
	return &dbHandler
}

func New() error {
	log.Println("validating DB details")
	if os.Getenv(DBType) == "" || os.Getenv(DBName) == "" || os.Getenv(DBHost) == "" || os.Getenv(DBPort) == "" ||
		os.Getenv(DBUser) == "" || os.Getenv(DBPassword) == "" {

		return fmt.Errorf("missing db details, please set below environment variables: "+
			"%v, %v, %v, %v, %v, %v\n", DBType, DBName, DBHost, DBPort, DBUser, DBPassword)
	}

	store := db{
		dbType:   os.Getenv(DBType),
		name:     os.Getenv(DBName),
		host:     os.Getenv(DBHost),
		port:     os.Getenv(DBPort),
		user:     os.Getenv(DBUser),
		password: os.Getenv(DBPassword),
	}

	var err error

	switch strings.ToLower(store.dbType) {
	case "postgres":
		dbHandler, err = newPostgresDB(store)
	default:
		return fmt.Errorf("storage type %v not supported, please check value of DB_TYPE in environment variables", store.dbType)
	}

	return err
}
