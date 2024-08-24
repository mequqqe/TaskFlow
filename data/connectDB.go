package db

import (
	"taskflow/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	var err error

	DBConnect := "host=localhost user=postgres password=1234 dbname=flowDB port=5432 sslmode=disable"

	// Open a connection to the PostgreSQL database using Gorm
	DB, err = gorm.Open(postgres.Open(DBConnect), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema for the Order and OrderStatusHistory models
	err = DB.AutoMigrate(&domain.Order{}, &domain.OrderStatusHistory{})
	if err != nil {
		return nil, err
	}

	return DB, nil
}
