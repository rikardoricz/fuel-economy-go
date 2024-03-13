package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-tracker/initializers"
	"github.com/rikardoricz/fuel-economy-tracker/models"
)

func CarsCreate(c *gin.Context) {
	var body struct {
		CarId        int
		CarName      string
		LicensePlate string
		CarType      string
		ProdDate     string
	}

	c.Bind(&body)

	car := models.Car{CarId: body.CarId, CarName: body.CarName, LicensePlate: body.LicensePlate, CarType: body.CarType, ProdDate: body.ProdDate}
	result := initializers.DB.Create(&car)

	if result.Error != nil {
		c.Status(400)
		return
	}

	c.JSON(200, gin.H{
		"car": car,
	})
}

func CarsIndex(c *gin.Context) {
	var cars []models.Car
	initializers.DB.Find(&cars)

	c.JSON(200, gin.H{
		"cars": cars,
	})
}

