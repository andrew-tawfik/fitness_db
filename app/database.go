package main

import (
	"fmt"
	"log"

	"fitness_db/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	dsn := "host=localhost user=postgres password=archSQL dbname=fitness_club port=5433 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")

	// Auto-migrate (creates tables)
	err = DB.AutoMigrate(
		&models.Member{},
		&models.Trainer{},
		&models.Class{},
		&models.HealthMetric{},
		&models.FitnessGoal{},
		&models.TrainingSession{},
		&models.ClassEnrollment{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")
}
