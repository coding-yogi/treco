package model

import (
	"gorm.io/gorm/clause"
	"regexp"
	"strings"
	"time"
	"treco/storage"
)

// Data
type Data struct {
	//DbHandler    *storage.DBHandler
	Jira         string
	ReportFormat string
	SuiteResult  SuiteResult
}

// SuiteResult
type SuiteResult struct {
	ID              uint    `gorm:"primarykey"`
	Build           string  `gorm:"uniqueIndex:ui_suite_result"`
	TestType        string  `gorm:"uniqueIndex:ui_suite_result"`
	Service         string  `gorm:"not null"`
	Environment     string  `gorm:"not null"`
	TimeTaken       float64 `gorm:"not null"`
	TotalExecuted   uint    `gorm:"default:0"`
	TotalPassed     uint    `gorm:"default:0"`
	TotalFailed     uint    `gorm:"default:0"`
	TotalSkipped    uint    `gorm:"default:0"`
	Coverage        float64 `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ScenarioResults []*ScenarioResult
}

// ScenarioResult
type ScenarioResult struct {
	ID            uint    `gorm:"primarykey"`
	ScenarioID    uint    `gorm:",not null"`
	SuiteResultID uint    `gorm:",not null"`
	Name          string  `gorm:"-"`
	Status        string  `gorm:",not null"`
	TimeTaken     float64 `gorm:"default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Scenario
type Scenario struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex:ui_scenario"`
	TestType  string    `gorm:"uniqueIndex:ui_scenario"`
	Service   string    `gorm:"uniqueIndex:ui_scenario"`
	Features  []Feature `gorm:"many2many:feature_scenarios"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Feature
type Feature struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Scenarios []Scenario `gorm:"many2many:feature_scenarios"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save
func (d *Data) Save() error {
	dbh := storage.Handler()

	suiteResult := &d.SuiteResult
	scenarioResults := suiteResult.ScenarioResults

	scenarios := make([]Scenario, 0, len(scenarioResults)) //scenarios

	// Loop through scenarios
	for _, scenarioResult := range scenarioResults {
		featureIds := getFeaturesFromScenario(d.Jira, scenarioResult.Name)
		features := make([]Feature, 0, len(featureIds)) //features

		for _, featureId := range featureIds {
			features = append(features, Feature{ID: featureId})
		}

		scenario := Scenario{
			Name:     scenarioResult.Name,
			TestType: d.SuiteResult.TestType,
			Service:  d.SuiteResult.Service,
			Features: features,
		}

		scenarios = append(scenarios, scenario)
	}

	switch db := (*dbh).(type) {
	case storage.Postgres:
		return writeToPostgres(&db, suiteResult, scenarios)
	}

	return nil
}

func writeToPostgres(db *storage.Postgres, suiteResult *SuiteResult, scenarios []Scenario) error {

	// Insert scenarios
	if err := db.GetDB().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}, {Name: "test_type"}, {Name: "service"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&scenarios).Error; err != nil {
		return err
	}

	// Update scenario results with scenario id
	for i, scenarioResult := range suiteResult.ScenarioResults {
		scenarioResult.ScenarioID = scenarios[i].ID
	}

	// Insert scenarios
	return db.GetDB().Create(suiteResult).Error
}

func getFeaturesFromScenario(p string, s string) []string {
	pat := `(?i)` + p + `-\d+`
	re := regexp.MustCompile(pat)
	matches := re.FindAllString(s, -1)

	for i := range matches {
		matches[i] = strings.ToUpper(matches[i])
	}

	return matches
}
