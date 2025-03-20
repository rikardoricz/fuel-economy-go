package main

import (
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func main() {
	err := initializers.DB.AutoMigrate(&models.Vehicle{})
	if err != nil {
		return
	}
	err = initializers.DB.AutoMigrate(&models.Refueling{})
	if err != nil {
		return
	}
}
