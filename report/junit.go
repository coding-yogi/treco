package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"treco/model"
)

// Expected Execution statuses
const (
	PASSED  = "passed"
	FAILED  = "failed"
	SKIPPED = "skipped"
)

// JunitReport struct
type JunitReport struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JunitTestSuite `xml:"testsuite"`
}

// JunitTestSuite struct
type JunitTestSuite struct {
	XMLName        xml.Name        `xml:"testsuite"`
	Tests          uint            `xml:"tests,attr"`
	Skipped        uint            `xml:"skipped,attr"`
	Failures       uint            `xml:"failures,attr"`
	Errors         uint            `xml:"errors,attr"`
	Time           float64         `xml:"time,attr"`
	JunitTestCases []JunitTestCase `xml:"testcase"`
}

// JunitTestCase struct
type JunitTestCase struct {
	XMLName  xml.Name  `xml:"testcase"`
	Name     string    `xml:"name,attr"`
	Class    string    `xml:"classname,attr"`
	Time     float64   `xml:"time,attr"`
	Features string    `xml:"features,attr"`
	Failure  *struct{} `xml:"failure,omitempty"`
	Skipped  *struct{} `xml:"skipped,omitempty"`
	Error    *struct{} `xml:"error,omitempty"`
}

var (
	errUnableToUnmarshalToJunit = "unmarshalling to junit failed"
)

type junitXMLParser struct{}

func (junitXMLParser) parse(r io.Reader, result *model.Data) error {
	suiteResult := &result.SuiteResult

	log.Println("reading report file")
	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("error reading: " + err.Error())
		return err
	}

	jur := JunitReport{}

	log.Println("unmarshalling to junit report")
	if strings.Contains(string(b), "testsuites") {
		if err = xml.Unmarshal(b, &jur); err != nil {
			return fmt.Errorf(errUnableToUnmarshalToJunit)
		}
	} else {
		suite := JunitTestSuite{}
		if err = xml.Unmarshal(b, &suite); err != nil {
			return fmt.Errorf(errUnableToUnmarshalToJunit)
		}

		jur.TestSuites = append(jur.TestSuites, suite)
	}

	for _, suite := range jur.TestSuites {
		suiteResult.TotalExecuted += suite.Tests
		suiteResult.TotalFailed += suite.Failures + suite.Errors
		suiteResult.TotalSkipped += suite.Skipped
		suiteResult.TotalPassed += suite.Tests - (suite.Failures + suite.Skipped + suite.Errors)
		suiteResult.TimeTaken += suite.Time

		for _, tc := range suite.JunitTestCases {
			status := PASSED
			if tc.Failure != nil || tc.Error != nil {
				status = FAILED
			} else if tc.Skipped != nil {
				status = SKIPPED
			}

			suiteResult.ScenarioResults = append(suiteResult.ScenarioResults, model.ScenarioResult{
				SuiteResultID: suiteResult.ID,
				Name:          tc.Name,
				Class:         tc.Class,
				Status:        status,
				TimeTaken:     tc.Time,
				Features:      strings.Split(tc.Features, " "),
			})
		}
	}

	return nil
}
