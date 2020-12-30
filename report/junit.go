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
	Tests          uint            `xml:"tests,attr"`
	Skipped        uint            `xml:"skipped,attr"`
	Failures       uint            `xml:"failures,attr"`
	Time           float64         `xml:"time,attr"`
	JunitTestCases []JunitTestCase `xml:"testcase"`
}

// JunitTestCase
type JunitTestCase struct {
	XMLName xml.Name  `xml:"testcase"`
	Name    string    `xml:"name,attr"`
	Time    float64   `xml:"time,attr"`
	Failure *struct{} `xml:"failure,omitempty"`
	Skipped *struct{} `xml:"skipped,omitempty"`
}

type junitXmlParser struct{}

func (junitXmlParser) parse(r io.Reader, result *model.Data) error {

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
			return err
		}
	} else {
		suite := JunitTestSuite{}
		if err = xml.Unmarshal(b, &suite); err != nil {
			return err
		}

		jur.TestSuites = append(jur.TestSuites, suite)
	}

	for _, s := range jur.TestSuites {
		suiteResult.TotalExecuted += s.Tests
		suiteResult.TotalFailed += s.Failures
		suiteResult.TotalSkipped += s.Skipped
		suiteResult.TotalPassed += s.Tests - (s.Failures + s.Skipped)
		suiteResult.TimeTaken += s.Time

		for _, u := range s.JunitTestCases {
			status := PASSED
			if u.Failure != nil {
				status = FAILED
			} else if u.Skipped != nil {
				status = SKIPPED
			}

			suiteResult.ScenarioResults = append(suiteResult.ScenarioResults, &model.ScenarioResult{
				SuiteResultID: suiteResult.ID,
				Name:          u.Name,
				Status:        status,
				TimeTaken:     u.Time,
			})
		}
	}

	return nil
}
