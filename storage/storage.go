package storage

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// DB Details
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
	DbType   string
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

// Handler returns the current DB handler
func Handler() *DBHandler {
	return &dbHandler
}

var (
	errMissingDBParams = fmt.Errorf("missing db details, please set below environment variables: "+
		"%v, %v, %v, %v, %v, %v\n", DBType, DBName, DBHost, DBPort, DBUser, DBPassword)

	errStrInvalidStorageType = "storage type %v not supported, please check value of DB_TYPE in environment variables"
)

// New initiates a new DB connection
func New() error {
	log.Println("validating DB details")
	if os.Getenv(DBType) == "" || os.Getenv(DBName) == "" || os.Getenv(DBHost) == "" || os.Getenv(DBPort) == "" ||
		os.Getenv(DBUser) == "" || os.Getenv(DBPassword) == "" {

		return errMissingDBParams
	}

	store := db{
		DbType:   os.Getenv(DBType),
		Name:     os.Getenv(DBName),
		Host:     os.Getenv(DBHost),
		Port:     os.Getenv(DBPort),
		User:     os.Getenv(DBUser),
		Password: os.Getenv(DBPassword),
	}

	var err error

	switch strings.ToLower(store.DbType) {
	case "postgres":
		dbHandler, err = newPostgresDB(store)
	default:
		return fmt.Errorf(errStrInvalidStorageType, store.DbType)
	}

	return err
}
