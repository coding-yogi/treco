package cmd

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestMissingFlags(t *testing.T) {

	cfg := config{
		Build:        "Test",
		Environment:  "dev",
		Jira:         "DAKOTA-000",
		ReportFile:   "../junit.xml",
		ReportFormat: "junit",
		Service:      "treco",
		TestType:     "unit",
	}

	configValue := reflect.ValueOf(&cfg).Elem()
	configType := configValue.Type()

	for fieldIdx := 0; fieldIdx < configValue.NumField(); fieldIdx++ {
		field := configValue.Field(fieldIdx)
		fieldName := configType.Field(fieldIdx).Name

		t.Run(fieldName, func(t *testing.T) {
			currentFieldValue := field.String()
			field.SetString("") //Set each field to empty

			err := validateFlags(cfg)
			require.Error(t, err, "no error returned for field "+fieldName)
			require.Equal(t, errMissingArguments, err)

			field.SetString(currentFieldValue) //Reset field
		})
	}
}

func TestInvalidHttpRequest(t *testing.T) {
	requestData := []struct {
		testName string
		request  *http.Request
		err      error
		status   int
	}{
		{
			testName: "method other that POST",
			request: &http.Request{
				Method: "GET",
				Header: map[string][]string{},
			},
			err:    fmt.Errorf(""),
			status: http.StatusMethodNotAllowed,
		},
		{
			testName: "content-type other than multipart/form-data",
			request: &http.Request{
				Method: "POST",
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			err:    fmt.Errorf("invalid content-type, expected: %v", expectedContentType),
			status: http.StatusBadRequest,
		},
		{
			testName: "missing request params",
			request: &http.Request{
				Method: "POST",
				Header: map[string][]string{
					"Content-Type": {expectedContentType},
				},
			},
			err:    fmt.Errorf("missing params: %v", strings.ToLower(strings.Join(requiredParams[:], ", "))),
			status: http.StatusBadRequest,
		},
	}

	for _, data := range requestData {
		t.Run(data.testName, func(t *testing.T) {
			status, err := validatePublishRequest(data.request)
			require.Error(t, err)
			require.Equal(t, data.err, err)
			require.Equal(t, data.status, status)
		})
	}
}

func TestValidHttpRequest(t *testing.T) {
	requestData := &http.Request{
		Method: "POST",
		Header: map[string][]string{
			"Content-Type": {expectedContentType},
		},
		Form: url.Values{
			strings.ToLower(BuildID):      {"test"},
			strings.ToLower(Environment):  {"dev"},
			strings.ToLower(Jira):         {"DAKOTA-007"},
			strings.ToLower(ReportFormat): {"junit"},
			strings.ToLower(Service):      {"test_service"},
			strings.ToLower(TestType):     {"unit"},
		},
	}

	status, err := validatePublishRequest(requestData)
	require.NoError(t, err)
	require.Equal(t, 0, status)
}

func TestInvalidParams(t *testing.T) {
	testData := []struct {
		testName   string
		testType   string
		reportType string
		err        error
	}{
		{
			testName:   "invalid test type",
			testType:   "unknown",
			reportType: "junit",
			err:        errInvalidTestType,
		},
		{
			testName:   "invalid report type",
			testType:   "unit",
			reportType: "mbunit",
			err:        errInvalidReportFormats,
		},
	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			err := validateParams(data.testType, data.reportType)
			require.Error(t, err)
			require.Equal(t, data.err, err)
		})
	}
}

func TestValidParams(t *testing.T) {
	for _, testType := range validTestTypes {
		err := validateParams(testType, validReportFormats[0])
		require.NoError(t, err)
	}
}
