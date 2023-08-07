package model

import (
	"testing"
	"treco/storage"

	"github.com/stretchr/testify/require"
)

func TestDataSave(t *testing.T) {
	data := &Data{
		Jira:         "Project",
		ReportFormat: "junit",
		SuiteResult: SuiteResult{
			TestType: "unit",
			Service:  "abc",
			ScenarioResults: []ScenarioResult{
				{
					Name:      "test-scenario-1 (project-123)",
					Status:    "passed",
					TimeTaken: 3.3,
				},
				{
					Name:      "(project-124)test-scenario-2(project-123)",
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
			projectName: "project",
			scenarioResult: ScenarioResult{
				Features: []string{"Project-123"},
			},
			featuresExtracted: []Feature{
				{
					ID: "PROJECT-123",
				},
			},
		}, {
			projectName: "Project",
			scenarioResult: ScenarioResult{
				Features: []string{"PROJECT-123", "project-456", "ProJecTIsGreat"},
			},
			featuresExtracted: []Feature{
				{
					ID: "PROJECT-123",
				},
				{
					ID: "PROJECT-456",
				},
			},
		}, {
			projectName: "project",
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
