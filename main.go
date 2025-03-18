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

func main() {
	r := gin.Default()

	r.POST("/vehicles", controllers.PostVehicles)
	r.GET("/vehicles", controllers.GetVehicles)
	r.GET("/vehicles/:id", controllers.GetVehicleByID)
	r.PUT("/vehicles/:id", controllers.UpdateVehicles)
	r.DELETE("/vehicles/:id", controllers.DeleteVehicles)

	err := r.Run()
	if err != nil {
		return
	}
}
