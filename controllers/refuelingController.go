package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"math"
	"net/http"
	"strconv"
	"time"
)

func GetRefuelings(c *gin.Context) {
	var refuelings []models.Refueling
	result := initializers.DB.Find(&refuelings)

	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
	}
	c.IndentedJSON(http.StatusOK, refuelings)
}

func CreateRefuelingForVehicle(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refueling models.Refueling
	if err := c.ShouldBindJSON(&refueling); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refueling.RefuelingDate = time.Now()
	refueling.VehicleID = uint(vehicleID)

	var previousRefueling models.Refueling
	result := initializers.DB.Where("vehicle_id = ?", vehicleID).Order("odometer_reading DESC").First(&previousRefueling)
	if result.Error == nil && previousRefueling.OdometerReading < refueling.OdometerReading {
		distanceTrip := refueling.OdometerReading - previousRefueling.OdometerReading
		fuelConsumption := (refueling.RefueledLiters / float64(distanceTrip)) * 100
		refueling.AvgFuelConsumption = math.Round(fuelConsumption*100) / 100
	}

	initializers.DB.Create(&refueling)
	c.IndentedJSON(http.StatusCreated, refueling)
}

func GetRefuelingsByVehicleID(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refueling []models.Refueling
	result := initializers.DB.Where("vehicle_id = ?", vehicleID).Find(&refueling)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, refueling)
}

func UpdateRefueling(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var body models.Refueling
	var refueling models.Refueling

	err = c.ShouldBindJSON(&refueling)
	if err != nil {
		return
	}

	initializers.DB.First(&refueling, vehicleID)
	initializers.DB.Model(&refueling).Updates(models.Refueling{
		RefuelingDate:      body.RefuelingDate.Truncate(24 * time.Hour),
		RefueledLiters:     body.RefueledLiters,
		OdometerReading:    body.OdometerReading,
		AvgFuelConsumption: body.AvgFuelConsumption})

	c.IndentedJSON(http.StatusOK, refueling)
}

func DeleteRefueling(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refueling models.Refueling

	result := initializers.DB.Where("vehicle_id = ?", vehicleID).Find(&refueling)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	initializers.DB.Delete(&refueling)
	c.IndentedJSON(http.StatusOK, refueling)
}

func GetRefuelingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refueling models.Refueling
	result := initializers.DB.First(&refueling, id)
	if result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, refueling)
}
