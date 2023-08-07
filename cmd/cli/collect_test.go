package cli

import (
	"reflect"
	"testing"
	"treco/conf"

	"github.com/stretchr/testify/require"
)

var testConfig = conf.Config{
	Build:        "Test",
	Environment:  "dev",
	Jira:         "DAKOTA",
	ReportFile:   "some_file.xml",
	ReportFormat: "junit",
	Service:      "treco",
	TestType:     "unit",
	Coverage:     "75.20",
}

func TestValidateFlagsWithMissingFlags(t *testing.T) {
	configValue := reflect.ValueOf(&testConfig).Elem()
	configType := configValue.Type()

	for fieldIdx := 0; fieldIdx < configValue.NumField(); fieldIdx++ {
		field := configValue.Field(fieldIdx)
		fieldName := configType.Field(fieldIdx).Name

		currentFieldValue := field.String()
		field.SetString("") //Set each field to empty

		t.Run(fieldName, func(t *testing.T) {
			err := validateFlags(testConfig)
			require.Error(t, err, "no error returned for field "+fieldName)
			require.Equal(t, errMissingArguments, err)
		})

		field.SetString(currentFieldValue) //Reset field
	}
}

func TestValidateFlagsWithValidFlags(t *testing.T) {
	err := validateFlags(testConfig)
	require.NoError(t, err)
}
