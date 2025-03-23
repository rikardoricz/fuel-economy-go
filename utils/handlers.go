package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func ParseID(c *gin.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, err
	}
	return id, nil
}

func CheckVehicleExists(id uint64) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := initializers.DB.First(&vehicle, id).Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func CheckRefuelingExists(id uint64) (*models.Refueling, error) {
	var refueling models.Refueling
	if err := initializers.DB.First(&refueling, id).Error; err != nil {
		return nil, err
	}
	return &refueling, nil
}

func HandleModelError(c *gin.Context, err error) bool {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return true
	}
	return false
}
