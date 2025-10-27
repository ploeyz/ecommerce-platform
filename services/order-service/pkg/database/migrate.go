package database

import (
	"log"

	"github.com/ploezy/ecommerce-platform/order-service/internal/models"
)

func AutoMigrate() error {
	log.Println("Running database migrations...")

	// Auto migrate models
	err := DB.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
	)	

	if err != nil{
		return err
	}
	log.Println("Database migration complete successfully")
	return nil
}