/*
Package storage handles persistent storage of report metrics
*/
package storage

import (
	"fmt"
	"log"
	"strings"
	"treco/conf"
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
	Setup(entities ...interface{}) error
}

type db struct {
	DBType   string
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
		"%v, %v, %v, %v, %v, %v", DBType, DBName, DBHost, DBPort, DBUser, DBPassword)

	errStrInvalidStorageType = "storage type %v not supported, please check value of DB_TYPE in environment variables"
)

// New initiates a new DB connection
func New() error {
	log.Println("validating DB details")
	if conf.Get(DBType) == "" || conf.Get(DBName) == "" || conf.Get(DBHost) == "" || conf.Get(DBPort) == "" ||
		conf.Get(DBUser) == "" || conf.Get(DBPassword) == "" {
		log.Printf("DBType: %v, DBName: %v, DBHost: %v, DBPort: %v, DBUser %v, DBPassword %v", conf.Get(DBType),
			conf.Get(DBName), conf.Get(DBHost), conf.Get(DBPort), conf.Get(DBUser), conf.Get(DBPassword))

		return errMissingDBParams
	}

	store := db{
		DBType:   conf.Get(DBType),
		Name:     conf.Get(DBName),
		Host:     conf.Get(DBHost),
		Port:     conf.Get(DBPort),
		User:     conf.Get(DBUser),
		Password: conf.Get(DBPassword),
	}

	var err error

	switch strings.ToLower(store.DBType) {
	case "postgres":
		dbHandler, err = newPostgresDB(store)
		return err

	default:
		return fmt.Errorf(errStrInvalidStorageType, store.DBType)
	}
}
