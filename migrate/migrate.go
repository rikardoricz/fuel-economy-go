package migrate

import (
	"github.com/rikardoricz/fuel-economy-tracker/initializers"
	"github.com/rikardoricz/fuel-economy-tracker/models"
)

func init() {
	initializers.LoadEnvVars()
	initializers.ConnectToDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.Car{})
}
