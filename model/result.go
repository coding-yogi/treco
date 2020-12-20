package model

import (
	"fmt"
	"log"
	"time"
	"treco/storage"
)

// Result
type Result struct {
	DbHandler     storage.DBHandler `pg:"-"`
	Build         string            `pg:",pk"`
	TestType      string
	Service       string
	TimeTaken     float64
	TotalExecuted int         `pg:",use_zero"`
	TotalPassed   int         `pg:",use_zero"`
	TotalFailed   int         `pg:",use_zero"`
	TotalSkipped  int         `pg:",use_zero"`
	Coverage      float64     `pg:",use_zero"`
	ExecutedAt    time.Time   `pg:"default:now()"`
	Scenarios     []*Scenario `pg:"rel:has-many"`
}

// Scenario
type Scenario struct {
	Id        int
	Build     string `pg:"fk:result_build"`
	Name      string
	Status    string
	TimeTaken float64 `pg:",use_zero"`
}

// Save result to storage
func (r *Result) Save() error {
	log.Println("saving report to database")

	switch dbt := r.DbHandler.(type) {
	case storage.Postgres:
		models := []interface{}{
			(*Result)(nil),
			(*Scenario)(nil),
		}

		log.Println("checking if schema exist")
		if err := dbt.Schema(models); err != nil {
			return err
		}

		log.Println("inserting in results DB")
		if err := dbt.Insert(r); err != nil {
			return err
		}

		log.Println("inserting in scenarios DB")
		if err := dbt.Insert(&r.Scenarios); err != nil {
			return err
		}

	default:
		return fmt.Errorf("non-supported db type %t", r.DbHandler)
	}

	return nil
}
