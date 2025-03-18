package models

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	LicensePlate   string
	Alias          string
	ProductionYear int
}
