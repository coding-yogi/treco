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

// JunitTestCase stuct
type JunitTestCase struct {
	XMLName  xml.Name  `xml:"testcase"`
	Name     string    `xml:"name,attr"`
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

	for _, s := range jur.TestSuites {
		suiteResult.TotalExecuted += s.Tests
		suiteResult.TotalFailed += s.Failures + s.Errors
		suiteResult.TotalSkipped += s.Skipped
		suiteResult.TotalPassed += s.Tests - (s.Failures + s.Skipped + s.Errors)
		suiteResult.TimeTaken += s.Time

		for _, u := range s.JunitTestCases {
			status := PASSED
			if u.Failure != nil || u.Error != nil {
				status = FAILED
			} else if u.Skipped != nil {
				status = SKIPPED
			}

			suiteResult.ScenarioResults = append(suiteResult.ScenarioResults, &model.ScenarioResult{
				SuiteResultID: suiteResult.ID,
				Name:          u.Name,
				Status:        status,
				TimeTaken:     u.Time,
				Features:      strings.Split(u.Features, " "),
			})
		}
	}

	return nil
}
