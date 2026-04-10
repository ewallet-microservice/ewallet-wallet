package bootstrap

import (
	"log"

	"github.com/mhasnanr/ewallet-wallet/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupDatabase() {
	dsn := GetEnv("CONNECTION_STRING", "")

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database")
	}

	Log.Infow("connected to database...")

	database.AutoMigrate(&models.Wallet{})
	DB = database
}
