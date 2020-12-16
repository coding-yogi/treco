package model

type Result struct {
	Service       string
	Build         string
	TestType      string
	Time          int64
	TotalExecuted int
	TotalPassed   int
	TotalFailed   int
	TotalSkipped  int
	Scenarios     []Scenario
}

type Scenario struct {
	Name   string
	Status string
	Time   float64
}
