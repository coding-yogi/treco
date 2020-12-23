package model

import (
	"fmt"
	"github.com/go-pg/pg/v10/orm"
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
	Id              int
	Build           string            `pg:",unique:suite_result"`
	TestType        string            `pg:",unique:suite_result"`
	Service         string            `pg:",notnull"`
	TimeTaken       float64           `pg:",notnull"`
	TotalExecuted   int               `pg:",use_zero"`
	TotalPassed     int               `pg:",use_zero"`
	TotalFailed     int               `pg:",use_zero"`
	TotalSkipped    int               `pg:",use_zero"`
	Coverage        float64           `pg:",use_zero"`
	ExecutedAt      time.Time         `pg:"default:now()"`
	ScenarioResults []*ScenarioResult `pg:"rel:has-many"`
}

// ScenarioResult
type ScenarioResult struct {
	Id            int
	ScenarioId    int          `pg:",notnull"`
	SuiteResultId int          `pg:",notnull"`
	Name          string       `pg:"-"`
	Status        string       `pg:",notnull"`
	TimeTaken     float64      `pg:",use_zero"`
	SuiteResult   *SuiteResult `pg:"rel:has-one"`
	Scenario      *Scenario    `pg:"rel:has-one"`
}

// Scenario
type Scenario struct {
	Id         int
	Name       string   `pg:",unique:scenario"`
	TestType   string   `pg:",unique:scenario"`
	Service    string   `pg:",unique:scenario"`
	FeatureIds []string `pg:",array"`
}

// Feature
type Feature struct {
	Id    string `pg:",pk"`
	Title string
}

// Save
func (d *Data) Save() error {
	dbh := storage.Handler()

	suiteResult := &d.SuiteResult
	scenarioResults := suiteResult.ScenarioResults

	scenarios := make([]*Scenario, 0)       //scenarios
	features := make([]Feature, 0)          //features
	featureMap := make(map[string]struct{}) //feature map with unique features

	// Loop through scenarios
	for _, scenarioResult := range scenarioResults {
		featureIds := getFeaturesFromScenario(d.Jira, scenarioResult.Name)
		for _, featureId := range featureIds {
			if _, ok := featureMap[featureId]; !ok {
				featureMap[featureId] = struct{}{} //throw away value
				features = append(features, Feature{Id: featureId})
			}
		}

		scenario := &Scenario{
			Name:       scenarioResult.Name,
			TestType:   d.SuiteResult.TestType,
			Service:    d.SuiteResult.Service,
			FeatureIds: featureIds,
		}

		scenarios = append(scenarios, scenario)
	}

	switch db := (*dbh).(type) {
	case storage.Postgres:
		return d.writeToPostgres(db, suiteResult, scenarioResults, scenarios, features)
	}

	return nil
}

func (d *Data) writeToPostgres(db storage.Postgres, suiteResult *SuiteResult,  scenarioResults []*ScenarioResult,
	scenarios []*Scenario, features []Feature) error {
	var query *orm.Query

	//Insert suite result
	query = db.GetDB().Model(suiteResult).Returning("id")
	if err := db.Insert(query); err != nil {
		return err
	}

	// Insert scenarios
	query = db.GetDB().
		Model(&scenarios).
		Returning("id").
		OnConflict("ON CONSTRAINT scenarios_name_test_type_service_key DO UPDATE").
		Set("feature_ids = EXCLUDED.feature_ids")

	if err := db.Insert(query); err != nil {
		return err
	}

	//Insert features
	if len(features) > 0 {
		query = db.GetDB().Model(&features).Returning("id").OnConflict("DO NOTHING")
		if err := db.Insert(query); err != nil {
			return err
		}
	}

	for i, scenarioResult := range scenarioResults {
		scenarioResult.ScenarioId = scenarios[i].Id
		scenarioResult.SuiteResultId = suiteResult.Id
	}

	//Insert scenario results
	query = db.GetDB().Model(&scenarioResults).Returning("id")
	if err := db.Insert(query); err != nil {
		return err
	}

	return nil
}

func getFeaturesFromScenario(p string, s string) []string {
	p = fmt.Sprintf(`(?i)%v-\d+`, p)
	re := regexp.MustCompile(p)
	matches := re.FindAllString(s, -1)

	for i := range matches {
		matches[i] = strings.ToUpper(matches[i])
	}

	return matches
}
