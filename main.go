package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-go/controllers"
	"github.com/rikardoricz/fuel-economy-go/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	vehicles := r.Group("/vehicles")
	{
		vehicles.GET("", controllers.GetVehicles)
		vehicles.GET("/:id", controllers.GetVehicleByID)
		vehicles.POST("", controllers.CreateVehicles)
		vehicles.PUT("/:id", controllers.UpdateVehicles)
		vehicles.DELETE("/:id", controllers.DeleteVehicles)

		vehicles.GET("/:id/refuelings", controllers.GetRefuelingsByVehicleID)
		vehicles.POST("/:id/refuelings", controllers.CreateRefuelingForVehicle)
	}

	refuelings := r.Group("/refuelings")
	{
		refuelings.GET("", controllers.GetRefuelings)
		refuelings.GET("/:id", controllers.GetRefuelingByID)
		refuelings.PUT("/:id", controllers.UpdateRefueling)
		refuelings.DELETE("/:id", controllers.DeleteRefueling)
	}

	return r
}

func main() {
	r := setupRouter()
	err := r.Run()
	if err != nil {
		return
	}
}
