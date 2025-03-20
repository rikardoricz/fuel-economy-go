package models

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	LicensePlate   string `json:"license_plate"`
	Alias          string `json:"alias"`
	ProductionYear int    `json:"production_year"`
}
