package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
	}

	for _, data := range dataSet {
		features := getFeaturesFromScenario(data.projectName, data.scenarioName)
		require.ElementsMatch(t, data.featuresExtracted, features)
	}
}
