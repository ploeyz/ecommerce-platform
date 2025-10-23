package database

import (
	"log"

	"github.com/ploezy/ecommerce-platform/product-service/internal/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting database migration...")
	
	err := db.AutoMigrate(
		&model.Product{},
	)
	if err != nil{
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully")
	return nil
}