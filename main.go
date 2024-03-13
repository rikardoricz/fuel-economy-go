package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rikardoricz/fuel-economy-tracker/controllers"
	"github.com/rikardoricz/fuel-economy-tracker/initializers"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {
	r := gin.Default()

	r.GET("/cars", controllers.CarsIndex)
	r.POST("/cars", controllers.CarsCreate)

	r.Run()
}
