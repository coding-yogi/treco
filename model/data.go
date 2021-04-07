/*
Package model defines DAO
*/
package model

import (
	"regexp"
	"strings"
	"time"
	"treco/storage"

	"gorm.io/gorm/clause"
)

// Data from report
type Data struct {
	//DbHandler    *storage.DBHandler
	Jira         string
	ReportFormat string
	SuiteResult  SuiteResult
}

// SuiteResult with execution summary
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

// ScenarioResult struct with execution details
type ScenarioResult struct {
	ID            uint    `gorm:"primarykey"`
	ScenarioID    uint    `gorm:",not null"`
	SuiteResultID uint    `gorm:",not null"`
	Name          string  `gorm:"-"`
	Status        string  `gorm:",not null"`
	TimeTaken     float64 `gorm:"default:0"`
	Features      []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Scenario struct with details of scenario
type Scenario struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex:ui_scenario"`
	TestType  string    `gorm:"uniqueIndex:ui_scenario"`
	Service   string    `gorm:"uniqueIndex:ui_scenario"`
	Features  []Feature `gorm:"many2many:feature_scenarios"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Feature struct for Jiras
type Feature struct {
	ID        string `gorm:"primaryKey"`
	Title     string
	Scenarios []Scenario `gorm:"many2many:feature_scenarios"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Save data to DB
func (d *Data) Save(dbh *storage.DBHandler) error {
	suiteResult := &d.SuiteResult
	scenarioResults := suiteResult.ScenarioResults

	scenarios := make([]Scenario, 0, len(scenarioResults)) //scenarios

	// Loop through scenarios
	for _, scenarioResult := range scenarioResults {
		scenarios = append(scenarios, Scenario{
			Name:     scenarioResult.Name,
			TestType: d.SuiteResult.TestType,
			Service:  d.SuiteResult.Service,
			Features: getFeaturesFromScenarioResult(d.Jira, *scenarioResult),
		})
	}

	return saveToDB(dbh, suiteResult, scenarios)
}

func saveToDB(dbh *storage.DBHandler, suiteResult *SuiteResult, scenarios []Scenario) error {
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

	// Insert suiteResults
	return db.GetDB().Create(suiteResult).Error
}

func getFeaturesFromScenarioName(projectName string, scenario string) []Feature {
	pat := `(?i)` + projectName + `-\d+`
	re := regexp.MustCompile(pat)
	matches := re.FindAllString(scenario, -1)

	features := make([]Feature, 0, len(matches))

	for i := range matches {
		features = append(features, Feature{ID: strings.ToUpper(matches[i])})
	}

	return features
}

func getFeaturesFromScenarioResult(projectName string, r ScenarioResult) []Feature {
	pat := `(?i)` + projectName + `-\d+`
	re := regexp.MustCompile(pat)
	features := make([]Feature, 0, len(r.Features))
	for _, f := range r.Features {
		f = strings.ToUpper(f)
		if re.MatchString(f) {
			features = append(features, Feature{ID: strings.ToUpper(f)})
		}
	}

	return features
}
