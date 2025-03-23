package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"github.com/rikardoricz/fuel-economy-go/utils"
	"gorm.io/gorm"
	"net/http"
)

func CreateVehicles(c *gin.Context) {
	var body models.Vehicle

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON body"})
		return
	}

	vehicle := models.Vehicle{
		LicensePlate:   body.LicensePlate,
		Alias:          body.Alias,
		ProductionYear: body.ProductionYear}

	result := initializers.DB.Create(&vehicle)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle already exists"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle creation failed"})
		return
	}

	c.IndentedJSON(http.StatusCreated, vehicle)
}

func GetVehicles(c *gin.Context) {
	var vehicles []models.Vehicle
	result := initializers.DB.Find(&vehicles)

	if utils.HandleModelError(c, result.Error) {
		return
	}

	if result.RowsAffected == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No vehicles found"})
		return
	}

	c.IndentedJSON(http.StatusOK, vehicles)
}

func GetVehicleByID(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	vehicle, err := utils.CheckVehicleExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	c.IndentedJSON(http.StatusOK, vehicle)
}

func DeleteVehicles(c *gin.Context) {
	id, err := utils.ParseID(c)
	if err != nil {
		return
	}

	vehicle, err := utils.CheckVehicleExists(id)
	if utils.HandleModelError(c, err) {
		return
	}

	result := initializers.DB.Delete(vehicle)
	if utils.HandleModelError(c, result.Error) {
		return
	}

	c.IndentedJSON(http.StatusOK, vehicle)
}

func UpdateVehicles(c *gin.Context) {
	id := c.Param("id")
	var body models.Vehicle
	var vehicle models.Vehicle

	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := initializers.DB.First(&vehicle, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result := initializers.DB.Model(&vehicle).Updates(models.Vehicle{
		LicensePlate:   body.LicensePlate,
		Alias:          body.Alias,
		ProductionYear: body.ProductionYear,
	}); result.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, vehicle)
}
