package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/controllers"
	"github.com/rikardoricz/fuel-economy-go/initializers"
	"github.com/rikardoricz/fuel-economy-go/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	db.AutoMigrate(&models.Vehicle{})
	initializers.DB = db
	return db
}

func createTestVehicle(t *testing.T, db *gorm.DB) models.Vehicle {
	vehicle := models.Vehicle{
		LicensePlate:   "DW12345",
		Alias:          "Test car",
		ProductionYear: 2023,
	}

	result := db.Create(&vehicle)
	if result.Error != nil {
		t.Fatalf("Failed to create test vehicle: %v", result.Error)
	}
	return vehicle
}

func TestCreateVehicles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	r := gin.Default()
	r.POST("/vehicles", controllers.CreateVehicles)

	vehicle := models.Vehicle{
		LicensePlate:   "DW12345",
		Alias:          "Test car",
		ProductionYear: 2023,
	}

	jsonData, err := json.Marshal(vehicle)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	req, _ := http.NewRequest("POST", "/vehicles", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Vehicle
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotZero(t, response.ID)
	assert.Equal(t, vehicle.LicensePlate, response.LicensePlate)
	assert.Equal(t, vehicle.Alias, response.Alias)
	assert.Equal(t, vehicle.ProductionYear, response.ProductionYear)
}

func TestGetVehicles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	createTestVehicle(t, db)
	createTestVehicle(t, db)

	r := gin.Default()
	r.GET("/vehicles", controllers.GetVehicles)

	req, _ := http.NewRequest("GET", "/vehicles", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var vehicles []models.Vehicle
	err := json.Unmarshal(w.Body.Bytes(), &vehicles)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Len(t, vehicles, 2)
}

func TestGetVehicleByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	vehicle := createTestVehicle(t, db)

	r := gin.Default()
	r.GET("/vehicles/:id", controllers.GetVehicleByID)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/vehicles/%d", vehicle.ID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Vehicle
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, vehicle.ID, response.ID)
	assert.Equal(t, vehicle.LicensePlate, response.LicensePlate)
	assert.Equal(t, vehicle.Alias, response.Alias)
	assert.Equal(t, vehicle.ProductionYear, response.ProductionYear)
}

func TestGetVehicleByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	r := gin.Default()
	r.GET("/vehicles/:id", controllers.GetVehicleByID)

	req, _ := http.NewRequest("GET", "/vehicles/999", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateVehicles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	vehicle := createTestVehicle(t, db)

	r := gin.Default()
	r.PUT("/vehicles/:id", controllers.UpdateVehicles)

	updatedVehicle := models.Vehicle{
		LicensePlate:   "WI11223",
		Alias:          "Updated Car",
		ProductionYear: 2024,
	}

	jsonData, err := json.Marshal(updatedVehicle)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	req, _ := http.NewRequest("PUT", fmt.Sprintf("/vehicles/%d", vehicle.ID), bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Vehicle
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, vehicle.ID, response.ID)
	assert.Equal(t, updatedVehicle.LicensePlate, response.LicensePlate)
	assert.Equal(t, updatedVehicle.Alias, response.Alias)
	assert.Equal(t, updatedVehicle.ProductionYear, response.ProductionYear)
}

func TestUpdateVehicles_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	r := gin.Default()
	r.PUT("/vehicles/:id", controllers.UpdateVehicles)

	vehicleData := models.Vehicle{
		LicensePlate:   "WI11223",
		Alias:          "Updated Car",
		ProductionYear: 2024,
	}

	jsonData, err := json.Marshal(vehicleData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	req, _ := http.NewRequest("PUT", "/vehicles/999", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteVehicles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	vehicle := createTestVehicle(t, db)

	r := gin.Default()
	r.DELETE("/vehicles/:id", controllers.DeleteVehicles)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/vehicles/%d", vehicle.ID), nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var count int64
	db.Model(&models.Vehicle{}).Where("id = ?", vehicle.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteVehicles_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	r := gin.Default()
	r.DELETE("/vehicles/:id", controllers.DeleteVehicles)

	req, _ := http.NewRequest("DELETE", "/vehicles/999", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
