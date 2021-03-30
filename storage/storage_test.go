package storage

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"os"
	"testing"
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
		_ = os.Setenv(dbParam, "test")
	}

	//test
	err := New()
	require.Error(t, err)
	require.Equal(t, fmt.Errorf(errStrInvalidStorageType, os.Getenv(DBType)), err)
}

func TestValidPostgresDB(t *testing.T) {
	//Set data
	_ = os.Setenv(DBHost, "localhost")
	_ = os.Setenv(DBUser, "some_user")
	_ = os.Setenv(DBPassword, "some_password")
	_ = os.Setenv(DBPort, "5432")
	_ = os.Setenv(DBName, "some_db")
	_ = os.Setenv(DBType, "postgres")

	//Mock
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
	_ = os.Setenv(DBHost, "localhost")
	_ = os.Setenv(DBUser, "some_user")
	_ = os.Setenv(DBPassword, "some_password")
	_ = os.Setenv(DBPort, "5432")
	_ = os.Setenv(DBName, "some_db")
	_ = os.Setenv(DBType, "postgres")

	//Test
	err := New()
	require.Error(t, err)
	require.Contains(t, err.Error(), "dial error")
	require.Equal(t, Postgres{}, dbHandler)
}
