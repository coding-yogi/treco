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
			ScenarioResults: []ScenarioResult{
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

func TestGetFeaturesFromScenarioResult(t *testing.T) {
	dataSet := []struct {
		projectName       string
		scenarioResult    ScenarioResult
		featuresExtracted []Feature
	}{
		{
			projectName: "dakota",
			scenarioResult: ScenarioResult{
				Features: []string{"Dakota-123"},
			},
			featuresExtracted: []Feature{
				{
					ID: "DAKOTA-123",
				},
			},
		}, {
			projectName: "Dakota",
			scenarioResult: ScenarioResult{
				Features: []string{"DAKOTA-123", "dakota-456", "DaKoTaIsGreat"},
			},
			featuresExtracted: []Feature{
				{
					ID: "DAKOTA-123",
				},
				{
					ID: "DAKOTA-456",
				},
			},
		}, {
			projectName: "dakota",
			scenarioResult: ScenarioResult{
				Features: []string{"abc", "123"},
			},
			featuresExtracted: []Feature{},
		},
	}

	for _, data := range dataSet {
		features := getFeaturesFromScenarioResult(data.projectName, data.scenarioResult)
		require.ElementsMatch(t, data.featuresExtracted, features)
	}
}
