package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"time"
)

func GetRefuelings(c *gin.Context) {
	var refuelings []models.Refueling
	result := initializers.DB.Find(&refuelings)

	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, refuelings)
}

func CreateRefuelingForVehicle(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vehicle models.Vehicle
	if err := initializers.DB.First(&vehicle, vehicleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if result.Error == nil {
		distanceTrip := refueling.OdometerReading - previousRefueling.OdometerReading
		if distanceTrip > 0 {
			fuelConsumption := (refueling.RefueledLiters / float64(distanceTrip)) * 100
			refueling.AvgFuelConsumption = math.Round(fuelConsumption*100) / 100
		}
	}

	if err := initializers.DB.Create(&refueling).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, refueling)
}

func GetRefuelingsByVehicleID(c *gin.Context) {
	vehicleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var vehicle models.Vehicle
	if err := initializers.DB.First(&vehicle, vehicleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Vehicle not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var refuelings []models.Refueling
	result := initializers.DB.Where("vehicle_id = ?", vehicleID).Find(&refuelings)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(refuelings) == 0 {
		c.IndentedJSON(http.StatusOK, []models.Refueling{})
		return
	}

	c.IndentedJSON(http.StatusOK, refuelings)
}

func UpdateRefueling(c *gin.Context) {
	refuelingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var refueling models.Refueling
	if err := initializers.DB.First(&refueling, refuelingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Refueling not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var body models.Refueling
	if err := c.ShouldBindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updateFields models.Refueling
	if !body.RefuelingDate.IsZero() {
		updateFields.RefuelingDate = body.RefuelingDate.Truncate(24 * time.Hour)
	}

	updateFields.RefueledLiters = body.RefueledLiters
	updateFields.OdometerReading = body.OdometerReading

	if body.OdometerReading > 0 || body.RefueledLiters > 0 {
		var previousRefueling models.Refueling
		result := initializers.DB.Where("vehicle_id = ? AND odometer_reading < ?",
			refueling.VehicleID, refueling.OdometerReading).
			Order("odometer_reading DESC").First(&previousRefueling)

		if result.Error == nil {
			distanceTrip := refueling.OdometerReading - previousRefueling.OdometerReading
			if distanceTrip > 0 {
				fuelConsumption := (refueling.RefueledLiters / float64(distanceTrip)) * 100
				updateFields.AvgFuelConsumption = math.Round(fuelConsumption*100) / 100
			}
		} else {
			updateFields.AvgFuelConsumption = body.AvgFuelConsumption
		}
	}

	if err := initializers.DB.Model(&refueling).Updates(updateFields).Error; err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	initializers.DB.First(&refueling, refuelingID)
	c.IndentedJSON(http.StatusOK, refueling)
}

func DeleteRefueling(c *gin.Context) {
	id := c.Param("id")

	var refueling models.Refueling
	if err := initializers.DB.First(&refueling, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Refueling not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := initializers.DB.Delete(&refueling)
	if result.Error != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Refueling deleted", "id": id})
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
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Refueling not found"})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, refueling)
}
