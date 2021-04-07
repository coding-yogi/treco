package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var testConfig = config{
	Build:        "Test",
	Environment:  "dev",
	Jira:         "DAKOTA",
	ReportFile:   "some_file.xml",
	ReportFormat: "junit",
	Service:      "treco",
	TestType:     "unit",
}

var testRequestParams = map[string]string{
	strings.ToLower(BuildID):      "test",
	strings.ToLower(Environment):  "dev",
	strings.ToLower(Jira):         "DAKOTA",
	strings.ToLower(ReportFormat): "junit",
	strings.ToLower(Service):      "test_service",
	strings.ToLower(TestType):     "unit",
}

var testFileContent = `
	<?xml version="1.0" encoding="UTF-8"?>
	<testsuite skipped="1" hostname="testHost" name="dakota.app.ui.tests.TestsForOnboarding" tests="5" failures="1" timestamp="2021-03-23T19:25:34 SGT" time="6.286" errors="1">
		<testcase name="test_with_error" time="6.286" classname="some.test.Class">
		<error type="org.openqa.selenium.WebDriverException" message="org.openqa.selenium.WebDriverException: An unknown server-side error occurred">
			<![CDATA[org.openqa.selenium.WebDriverException: org.openqa.selenium.WebDriverException: An unknown server-side error occurred while processing the command. Original error: Cannot rewrite element locator 'get_started_button' to its complete form, because the current application package name is unknown. Consider providing the app package name or changing the locator to '<package_name>:id/get_started_button' format.
			]]>
		</error>
		</testcase> 
		<system-out/>
		<testcase name="test_skipped" time="0.0" classname="some.test.Class">
			<skipped/>
		</testcase> 
		<testcase name="test_failed" time="2.123" classname="some.test.Class">
			<failed/>
		</testcase> 
		<testcase name="test_passed_1" time="1.987" classname="some.test.Class"/>
		<testcase name="test_passed_2" time="3.14" classname="some.test.Class"/>
	</testsuite> 
	`

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

//nolint: scopelint
func TestValidateParamsWithInvalidParams(t *testing.T) {
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

func TestValidateParamsValidParams(t *testing.T) {
	for _, testType := range validTestTypes {
		err := validateParams(testType, validReportFormats[0])
		require.NoError(t, err)
	}
}

func TestErrorResponse(t *testing.T) {
	resRecorder := httptest.NewRecorder()
	errDescription := "some error occured"
	errCode := 400
	err := fmt.Errorf(errDescription)

	errorResponse := Error{
		Code:        errCode,
		Description: errDescription,
	}

	body, _ := json.Marshal(errorResponse)

	sendErrorResponse(resRecorder, err, errDescription, errCode)
	require.Equal(t, errCode, resRecorder.Code)
	require.Equal(t, "application/json", resRecorder.Header().Get("content-type"))
	require.Equal(t, body, resRecorder.Body.Bytes())
}

//nolint: scopelint
func TestPublishHandlerWithInvalidData(t *testing.T) {
	type testData struct {
		testName string
		request  *http.Request
		resErr   Error
	}

	type feederFunc func() testData

	feeder := []feederFunc{
		func() testData {
			request, err := createTestHTTPRequest("GET", "multipart/form-data",
				testRequestParams, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "method other that POST",
				request:  request,
				resErr: Error{
					Code:        http.StatusMethodNotAllowed,
					Description: "",
				},
			}
		},
		func() testData {
			request, err := createTestHTTPRequest("POST", "application/json",
				testRequestParams, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "content-type other than multipart/form-data",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: fmt.Sprintf("invalid content-type, expected: %v", expectedContentType),
				},
			}
		},
		func() testData {
			request, err := createTestHTTPRequest("POST", "multipart/form-data",
				map[string]string{}, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "missing request params",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: fmt.Sprintf("missing params: %v", strings.ToLower(strings.Join(requiredParams[:], ", "))),
				},
			}
		},
		func() testData {
			requestParams := make(map[string]string)
			for k, v := range testRequestParams {
				requestParams[k] = v
			}

			requestParams[strings.ToLower(TestType)] = "unknown"
			request, err := createTestHTTPRequest("POST", "multipart/form-data",
				requestParams, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "invalid params",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: errInvalidTestType.Error(),
				},
			}
		},
		func() testData {
			request, err := createTestHTTPRequest("POST", "multipart/form-data",
				testRequestParams, "some invalid text format")
			require.NoError(t, err)

			return testData{
				testName: "invalid file contents",
				request:  request,
				resErr: Error{
					Code:        http.StatusInternalServerError,
					Description: "unable to process the request",
				},
			}
		},
	}

	for _, feederFun := range feeder {
		data := feederFun()
		t.Run(data.testName, func(t *testing.T) {
			resRecorder := httptest.NewRecorder()
			publishHandler(resRecorder, data.request)

			body, _ := json.Marshal(data.resErr)

			require.Equal(t, data.resErr.Code, resRecorder.Code)
			require.Equal(t, "application/json", resRecorder.Header().Get("content-type"))
			require.Equal(t, string(body), resRecorder.Body.String())
		})
	}
}

func TestPublishHandlerWithValidRequest(t *testing.T) {
	req, err := createTestHTTPRequest("POST", "multipart/form-data", testRequestParams, testFileContent)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	publishHandler(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}

func createTestHTTPRequest(method, contentType string, params map[string]string, fileContents string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//Write text params
	for k, v := range params {
		err := writer.WriteField(k, v)
		if err != nil {
			return &http.Request{}, err
		}
	}

	//Add file to multipart data
	part, _ := writer.CreateFormFile(strings.ToLower(ReportFile), "some_file.xml")
	_, err := io.Copy(part, bytes.NewReader([]byte(fileContents)))
	if err != nil {
		return &http.Request{}, err
	}

	//Close writer
	err = writer.Close()
	if err != nil {
		return &http.Request{}, err
	}

	//Create request
	req := httptest.NewRequest(method, "/treco/v1/publish/report", body)
	if contentType == "multipart/form-data" {
		req.Header.Set("Content-Type", writer.FormDataContentType())
		err = req.ParseMultipartForm(10 << 20)
		if err != nil {
			return &http.Request{}, err
		}
	} else {
		req.Header.Set("Content-Type", contentType)
	}

	return req, err
}
