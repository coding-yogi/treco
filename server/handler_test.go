package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	ContentTypeHeader            = "content-type"
	ContentTypeMultipartFormData = "multipart/form-data"
	ContentTypeApplicationJSON   = "application/json"
	MethodPost                   = "POST"
	MethodGet                    = "GET"
)

var testRequestParams = map[string]string{
	strings.ToLower(BuildID):      "test",
	strings.ToLower(Environment):  "dev",
	strings.ToLower(Jira):         "DAKOTA",
	strings.ToLower(ReportFormat): "junit",
	strings.ToLower(Service):      "test_service",
	strings.ToLower(TestType):     "unit",
	strings.ToLower(Coverage):     "10.0",
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

// nolint: scopelint
func TestValidateParamsWithInvalidParams(t *testing.T) {
	testData := []struct {
		testName   string
		testType   string
		reportType string
		coverage   string
		err        error
	}{
		{
			testName:   "invalid test type",
			testType:   "unknown",
			reportType: "junit",
			coverage:   "0.0",
			err:        fmt.Errorf(errInvalidTestType, "unknown", validTestTypes),
		},
		{
			testName:   "invalid report type",
			testType:   "unit",
			reportType: "mbunit",
			coverage:   "20.10",
			err:        fmt.Errorf(errInvalidReportFormats, "mbunit", validReportFormats),
		},
		{
			testName:   "invalid coverage",
			testType:   "unit",
			reportType: "junit",
			coverage:   "abc",
			err:        fmt.Errorf(errCoverageValueNotFloat),
		},
	}

	for _, data := range testData {
		t.Run(data.testName, func(t *testing.T) {
			err := ValidateParams(data.testType, data.reportType, data.coverage)
			require.Error(t, err)
			require.Equal(t, data.err, err)
		})
	}
}

func TestValidateParamsWithValidValues(t *testing.T) {
	for _, testType := range validTestTypes {
		err := ValidateParams(testType, validReportFormats[0], "0.0")
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
	require.Equal(t, ContentTypeApplicationJSON, resRecorder.Header().Get(ContentTypeHeader))
	require.Equal(t, body, resRecorder.Body.Bytes())
}

// nolint: scopelint
func TestPublishHandlerWithInvalidData(t *testing.T) {
	type testData struct {
		testName string
		request  *http.Request
		resErr   Error
	}

	type feederFunc func() testData

	feeder := []feederFunc{
		func() testData {
			request, err := createTestHTTPRequest(MethodGet, ContentTypeMultipartFormData,
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
			request, err := createTestHTTPRequest(MethodPost, ContentTypeApplicationJSON,
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
			request, err := createTestHTTPRequest(MethodPost, ContentTypeMultipartFormData,
				map[string]string{}, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "missing request params",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: fmt.Sprintf("missing params: %v", strings.ToLower(strings.Join(RequiredParams[:], ", "))),
				},
			}
		},
		func() testData {
			requestParams := make(map[string]string)
			for k, v := range testRequestParams {
				requestParams[k] = v
			}

			invalidTestType := "unknown"
			requestParams[strings.ToLower(TestType)] = invalidTestType
			request, err := createTestHTTPRequest(MethodPost, ContentTypeMultipartFormData,
				requestParams, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "invalid params",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: fmt.Errorf(errInvalidTestType, invalidTestType, validTestTypes).Error(),
				},
			}
		},
		func() testData {
			request, err := createTestHTTPRequest(MethodPost, ContentTypeMultipartFormData,
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
		func() testData {
			requestParams := make(map[string]string)
			for k, v := range testRequestParams {
				requestParams[k] = v
			}

			requestParams[strings.ToLower(Coverage)] = "unknown"
			request, err := createTestHTTPRequest(MethodPost, ContentTypeMultipartFormData,
				requestParams, testFileContent)
			require.NoError(t, err)

			return testData{
				testName: "invalid coverage value",
				request:  request,
				resErr: Error{
					Code:        http.StatusBadRequest,
					Description: fmt.Errorf(errCoverageValueNotFloat).Error(),
				},
			}
		},
	}

	for _, feederFun := range feeder {
		data := feederFun()
		t.Run(data.testName, func(t *testing.T) {
			resRecorder := httptest.NewRecorder()

			publishHandler := PublishHandler{}
			publishHandler.ServeHTTP(resRecorder, data.request)

			body, _ := json.Marshal(data.resErr)

			require.Equal(t, data.resErr.Code, resRecorder.Code)
			require.Equal(t, ContentTypeApplicationJSON, resRecorder.Header().Get(ContentTypeHeader))
			require.Equal(t, string(body), resRecorder.Body.String())
		})
	}
}

func TestPublishHandlerWithValidRequest(t *testing.T) {
	req, err := createTestHTTPRequest(MethodPost, ContentTypeMultipartFormData, testRequestParams, testFileContent)
	require.NoError(t, err)

	res := httptest.NewRecorder()
	publishHandler := PublishHandler{}
	publishHandler.ServeHTTP(res, req)

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
	req := httptest.NewRequest(method, "/v1/publish/report", body)
	if contentType == ContentTypeMultipartFormData {
		req.Header.Set(ContentTypeHeader, writer.FormDataContentType())
		err = req.ParseMultipartForm(10 << 20)
		if err != nil {
			return &http.Request{}, err
		}
	} else {
		req.Header.Set(ContentTypeHeader, contentType)
	}

	return req, err
}
