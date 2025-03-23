package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"github.com/rikardoricz/fuel-economy-go/utils"
	"math"
	"net/http"
	"time"
)

func CreateRefuelingForVehicle(c *gin.Context) {
	vehicleID, err := utils.ParseID(c)
	if err != nil {
		return
	}

	_, err = utils.CheckVehicleExists(vehicleID)
	if utils.HandleModelError(c, err) {
		return
	}

	var refueling models.Refueling
	if err := c.ShouldBindJSON(&refueling); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refueling.VehicleID = uint(vehicleID)

	calculateFuelConsumption(&refueling)

	if err := initializers.DB.Create(&refueling).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, refueling)
}

func GetRefuelings(c *gin.Context) {
	var refuelings []models.Refueling
	result := initializers.DB.Find(&refuelings)

	if utils.HandleModelError(c, result.Error) {
		return
	}

	c.IndentedJSON(http.StatusOK, refuelings)
}

func GetRefuelingByID(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	refueling, err := utils.CheckRefuelingExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	c.IndentedJSON(http.StatusOK, refueling)
}

func GetRefuelingsByVehicleID(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	_, err = utils.CheckVehicleExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	var refuelings []models.Refueling
	result := initializers.DB.Where("vehicle_id = ?", id).Find(&refuelings)
	if utils.HandleModelError(c, result.Error) {
		return
	}

	c.IndentedJSON(http.StatusOK, refuelings)
}

func UpdateRefueling(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	refueling, err := utils.CheckRefuelingExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	var body models.Refueling
	if err := c.ShouldBindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updateFields := prepareRefuelingUpdate(refueling, body)

	if err := initializers.DB.Model(refueling).Updates(updateFields).Error; err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	initializers.DB.First(refueling, id)
	c.IndentedJSON(http.StatusOK, refueling)
}

func DeleteRefueling(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	refueling, err := utils.CheckRefuelingExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	result := initializers.DB.Delete(refueling)
	if utils.HandleModelError(c, result.Error) {
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Refueling deleted", "id": id})
}

func calculateFuelConsumption(refueling *models.Refueling) {
	var previousRefueling models.Refueling
	result := initializers.DB.Where("vehicle_id = ?", refueling.VehicleID).
		Order("odometer_reading DESC").First(&previousRefueling)

	if result.Error == nil {
		distanceTrip := refueling.OdometerReading - previousRefueling.OdometerReading
		if distanceTrip > 0 {
			fuelConsumption := (refueling.RefueledLiters / float64(distanceTrip)) * 100
			refueling.AvgFuelConsumption = math.Round(fuelConsumption*100) / 100
		}
	}
}

func prepareRefuelingUpdate(original *models.Refueling, updates models.Refueling) models.Refueling {
	var updateFields models.Refueling

	if !updates.RefuelingDate.IsZero() {
		updateFields.RefuelingDate = updates.RefuelingDate.Truncate(24 * time.Hour)
	}

	updateFields.RefueledLiters = updates.RefueledLiters
	updateFields.OdometerReading = updates.OdometerReading

	if updates.OdometerReading > 0 || updates.RefueledLiters > 0 {
		var previousRefueling models.Refueling
		result := initializers.DB.Where("vehicle_id = ? AND odometer_reading < ?",
			original.VehicleID, original.OdometerReading).
			Order("odometer_reading DESC").First(&previousRefueling)

		if result.Error == nil {
			distanceTrip := original.OdometerReading - previousRefueling.OdometerReading
			if distanceTrip > 0 {
				fuelConsumption := (original.RefueledLiters / float64(distanceTrip)) * 100
				updateFields.AvgFuelConsumption = math.Round(fuelConsumption*100) / 100
			}
		} else {
			updateFields.AvgFuelConsumption = updates.AvgFuelConsumption
		}
	}

	return updateFields
}
