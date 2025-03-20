package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"net/http"
)

func CreateVehicles(c *gin.Context) {
	var body models.Vehicle

	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	vehicle := models.Vehicle{LicensePlate: body.LicensePlate, Alias: body.Alias, ProductionYear: body.ProductionYear}
	result := initializers.DB.Create(&vehicle)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	}
	c.IndentedJSON(http.StatusCreated, vehicle)
}

func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle

	result := initializers.DB.Find(&vehicles)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	}
	c.IndentedJSON(http.StatusOK, vehicles)
}

func GetVehicleByID(c *gin.Context) {
	var vehicle models.Vehicle
	id := c.Param("id")

	initializers.DB.First(&vehicle, id)

	c.IndentedJSON(http.StatusOK, vehicle)
}

func DeleteVehicles(c *gin.Context) {
	id := c.Param("id")
	result := initializers.DB.Delete(&models.Vehicle{}, id)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Vehicle deleted", "id": id})
}

func UpdateVehicles(c *gin.Context) {
	id := c.Param("id")
	var body models.Vehicle
	var vehicle models.Vehicle

	err := c.BindJSON(&body)
	if err != nil {
		return
	}

	initializers.DB.First(&vehicle, id)
	initializers.DB.Model(&vehicle).Updates(models.Vehicle{
		LicensePlate:   body.LicensePlate,
		Alias:          body.Alias,
		ProductionYear: body.ProductionYear})

	c.IndentedJSON(http.StatusOK, vehicle)
}
