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

const (
	PASSED  = "passed"
	FAILED  = "failed"
	SKIPPED = "skipped"
)

// JunitReport
type JunitReport struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JunitTestSuite `xml:"testsuite"`
}

// JunitTestSuite
type JunitTestSuite struct {
	XMLName        xml.Name        `xml:"testsuite"`
	Tests          int             `xml:"tests,attr"`
	Skipped        int             `xml:"skipped,attr"`
	Failures       int             `xml:"failures,attr"`
	Time           float64         `xml:"time,attr"`
	JunitTestCases []JunitTestCase `xml:"testcase"`
}

// JunitTestCase
type JunitTestCase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Time    float64  `xml:"time,attr"`
	Failure Failure  `xml:"failure,omitempty"`
	Skipped Skipped  `xml:"skipped,omitempty"`
}

// Failure
type Failure struct {
	Message string `xml:"message,attr"`
}

// Skipped
type Skipped struct {}

type junitXmlParser struct{}

func (junitXmlParser) parse(r io.Reader, result *model.Result) error {
	log.Println("reading report file")
	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("error reading: " + err.Error())
		return err
	}

	jur := &JunitReport{}

	log.Println("unmarshalling to junit report")
	if strings.Contains(string(b), "testsuites") {
		if err = xml.Unmarshal(b, jur); err != nil {
			return err
		}
	} else {
		suite := &JunitTestSuite{}
		if err = xml.Unmarshal(b, suite); err != nil {
			return err
		}

		jur.TestSuites = append(jur.TestSuites, *suite)
	}

	for _, s := range jur.TestSuites {
		result.TotalExecuted = +s.Tests
		result.TotalFailed = +s.Failures
		result.TotalSkipped = +s.Skipped
		result.TotalPassed = +(s.Tests - (s.Failures + s.Skipped))
		result.TimeTaken = +s.Time

		for _, u := range s.JunitTestCases {
			name := u.Name
			status := PASSED
			if strings.ToLower(u.Failure.Message) != "" {
				status = FAILED
			} else if u.Time == 0.0 {
				status = SKIPPED
			}

			time := u.Time
			result.Scenarios = append(result.Scenarios, &model.Scenario{Build: result.Build, Name: name, Status: status, TimeTaken: time})
		}
	}

	return nil
}
