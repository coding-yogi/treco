package model

import (
	"testing"
	"treco/storage"

	"github.com/stretchr/testify/require"
)

func TestDataSave(t *testing.T) {
	data := &Data{
		Jira:         "Dakota",
		ReportFormat: "junit",
		SuiteResult: SuiteResult{
			TestType: "unit",
			Service:  "abc",
			ScenarioResults: []*ScenarioResult{
				{
					Name:      "test-scenario-1 (dakota-123)",
					Status:    "passed",
					TimeTaken: 3.3,
				},
				{
					Name:      "(dakota-124)test-scenario-2(dakota-123)",
					Status:    "failed",
					TimeTaken: 1.3,
				},
			},
		},
	}

	err := data.Save(storage.Handler())
	require.NoError(t, err)
}

func TestGetFeaturesFromScenario(t *testing.T) {
	dataSet := []struct {
		projectName       string
		scenarioName      string
		featuresExtracted []string
	}{
		{
			projectName:       "dakota",
			scenarioName:      "some test (Dakota-123)",
			featuresExtracted: []string{"DAKOTA-123"},
		}, {
			projectName:       "Dakota",
			scenarioName:      "Dakota-456 some test (Dakota-123)",
			featuresExtracted: []string{"DAKOTA-456", "DAKOTA-123"},
		}, {
			projectName:       "DAKOTA",
			scenarioName:      "Dakota-456 some test (Dakota-123)",
			featuresExtracted: []string{"DAKOTA-456", "DAKOTA-123"},
		},
		{
			projectName:       "DAKOTA",
			scenarioName:      "some test dakota",
			featuresExtracted: []string{},
		},
		{
			projectName:       "DAKOTA",
			scenarioName:      "Dakota-456 some test Dakota-456 dakota",
			featuresExtracted: []string{"DAKOTA-456", "DAKOTA-456"},
		},
	}

	for _, data := range dataSet {
		features := getFeaturesFromScenario(data.projectName, data.scenarioName)
		require.ElementsMatch(t, data.featuresExtracted, features)
	}
}
