package models

type FuelType int

const (
	Euro95E10 FuelType = iota
	Euro95E5
	Diesel
)

func (t FuelType) String() string {
	return [...]string{"Euro95 E10", "Euro95 E5", "Diesel"}[t]
}