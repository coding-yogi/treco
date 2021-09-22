package storage

import (
	"fmt"
	"os"
	"testing"
	"treco/conf"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var dbParams = []string{DBHost, DBPort, DBUser, DBPassword, DBName, DBType}

func TestMissingDBDetails(t *testing.T) {
	//Set some value to all db params
	for _, dbParam := range dbParams {
		_ = os.Setenv(dbParam, "test")
	}

	for _, dbParam := range dbParams {
		_ = os.Unsetenv(dbParam) //unset the value

		//run test
		t.Run(dbParam, func(t *testing.T) {
			err := New()
			require.Error(t, err)
			require.Equal(t, errMissingDBParams, err)
		})

		_ = os.Setenv(dbParam, "test") //reset the value
	}
}

func TestInvalidDBType(t *testing.T) {
	//Set some value to all db params
	for _, dbParam := range dbParams {
		conf.Set(dbParam, "test")
	}

	//test
	err := New()
	require.Error(t, err)
	require.Equal(t, fmt.Errorf(errStrInvalidStorageType, os.Getenv(DBType)), err)
}

func TestValidPostgresDBType(t *testing.T) {
	//Set data
	conf.Set(DBHost, "localhost")
	conf.Set(DBUser, "some_user")
	conf.Set(DBPassword, "some_password")
	conf.Set(DBPort, "5432")
	conf.Set(DBName, "some_db")
	conf.Set(DBType, "postgres")

	//Mock
	connectToDB := connectToPostgresDB
	connectToPostgresDB = func(dsn string) (*gorm.DB, error) {
		return &gorm.DB{}, nil
	}

	t.Run("valid db", func(t *testing.T) {
		//Test
		err := New()
		require.NoError(t, err)
		require.IsType(t, Postgres{}, dbHandler)
	})

	//Reset
	connectToPostgresDB = connectToDB
}

func TestInvalidPostgresDBConnection(t *testing.T) {
	//Set data
	conf.Set(DBHost, "localhost")
	conf.Set(DBUser, "some_user")
	conf.Set(DBPassword, "some_password")
	conf.Set(DBPort, "5432")
	conf.Set(DBName, "some_db")
	conf.Set(DBType, "postgres")

	//Test
	err := New()
	require.Error(t, err)
	require.Contains(t, err.Error(), "dial error")
	require.Equal(t, Postgres{}, dbHandler)
}

func TestGetHandler(t *testing.T) {
	handler := Handler()
	require.IsType(t, (*DBHandler)(nil), handler)
}
