package models

type Trainer struct {
    TrainerID      uint   `gorm:"primaryKey;autoIncrement"`
    FirstName      string `gorm:"size:50;not null"`
    LastName       string `gorm:"size:50;not null"`
    Email          string `gorm:"size:100;uniqueIndex;not null"`
    Specialization string `gorm:"size:100"`
    Phone          string `gorm:"size:15"`
    
    // Relationships
    Classes          []Class           `gorm:"foreignKey:TrainerID"`
    TrainingSessions []TrainingSession `gorm:"foreignKey:TrainerID"`
}