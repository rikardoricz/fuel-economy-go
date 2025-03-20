package models

import (
	"gorm.io/gorm"
	"time"
)

type Refueling struct {
	gorm.Model
	RefuelingDate      time.Time `json:"refueling_date"`
	RefueledLiters     float64   `json:"refueled_liters"`
	OdometerReading    int       `json:"odometer_reading"`
	AvgFuelConsumption float64   `json:"avg_fuel_consumption"`
	VehicleID          uint      `json:"vehicle_id"`
}
