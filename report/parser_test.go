package report

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"treco/model"
)

func TestInvalidReportFormat(t *testing.T) {
	reportFormat := "unsupported format"
	data := &model.Data{
		ReportFormat: reportFormat,
	}

	contents := "test"
	err := Parse(bytes.NewReader([]byte(contents)), data)
	require.Equal(t, fmt.Errorf(errInvalidReportType, reportFormat), err)
}

func TestInvalidJunitContent(t *testing.T) {
	data := &model.Data{
		ReportFormat: "junit",
	}

	contents := "test"
	err := Parse(bytes.NewReader([]byte(contents)), data)
	require.Equal(t, fmt.Errorf(errUnableToUnmarshalToJunit), err)
}

func TestJunitReportParsing(t *testing.T) {
	data := &model.Data{
		ReportFormat: "junit",
	}

	contents := `
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
	err := Parse(bytes.NewReader([]byte(contents)), data)
	require.NoError(t, err, "Paring error")
	require.Equal(t, data.SuiteResult.TotalExecuted, uint(5))
	require.Equal(t, data.SuiteResult.TotalFailed, uint(2))
	require.Equal(t, data.SuiteResult.TotalSkipped, uint(1))
	require.Equal(t, data.SuiteResult.TotalPassed, uint(2))
	require.Equal(t, len(data.SuiteResult.ScenarioResults), 5)
}
