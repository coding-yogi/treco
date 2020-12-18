package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"treco/model"
)

const (
	PASSED  = "passed"
	FAILED  = "failed"
	SKIPPED = "skipped"
)

type JunitReport struct {
	XMLName    xml.Name         `xml:"testsuites"`
	TestSuites []JunitTestSuite `xml:"testsuite"`
}

type JunitTestSuite struct {
	XMLName        xml.Name        `xml:"testsuite"`
	Tests          int             `xml:"tests,attr"`
	Skipped        int             `xml:"skipped,attr"`
	Failures       int             `xml:"failures,attr"`
	Time           float64         `xml:"time,attr"`
	JunitTestCases []JunitTestCase `xml:"testcase"`
}

type JunitTestCase struct {
	XMLName xml.Name `xml:"testcase"`
	Name    string   `xml:"name,attr"`
	Time    float64  `xml:"time,attr"`
	Failure Failure  `xml:"failure,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
}

type JunitXmlParser struct{}

func (JunitXmlParser) Parse(r io.Reader, result *model.Result) error {

	b, err := ioutil.ReadAll(r)
	if err != nil {
		fmt.Println("error reading: " + err.Error())
		return err
	}

	jur := &JunitReport{}

	if strings.Contains(string(b), "testsuites") {
		err = xml.Unmarshal(b, jur)
		if err != nil {
			fmt.Println("error unmarshalling junit report: " + err.Error())
			return err
		}
	} else {
		suite := &JunitTestSuite{}
		err = xml.Unmarshal(b, suite)
		if err != nil {
			fmt.Println("error unmarshalling junit report: " + err.Error())
			return err
		}

		jur.TestSuites = append(jur.TestSuites, *suite)
	}

	for _, s := range jur.TestSuites {
		result.TotalExecuted = +s.Tests
		result.TotalFailed = +s.Failures
		result.TotalSkipped = +s.Skipped
		result.TotalPassed = +(s.Tests - (s.Failures + s.Skipped))

		for _, u := range s.JunitTestCases {
			name := u.Name
			status := PASSED
			if strings.ToLower(u.Failure.Message) != "" {
				status = FAILED
			}
			time := u.Time
			result.Scenarios = append(result.Scenarios, model.Scenario{Name: name, Status: status, Time: time})
		}
	}

	fmt.Println(result)
	return nil
}
