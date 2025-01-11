package types

import "time"

type EnergyType int64

const (
	Electricity EnergyType = iota
	Gas
	Water
)

type EnergyReading struct {
	Time  time.Time  `json:"time"`
	Type  EnergyType `json:"type"`
	Unit  string     `json:"unit"`
	Value float64    `json:"value"`
	Error error      `json:"error"`
}

type MultipleEnergyReadings struct {
	EnergyReadings []EnergyReading `json:"energy_readings"`
	Unit           string          `json:"unit"`
	Error          error           `json:"error"`
}
