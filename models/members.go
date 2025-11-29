package models

import (
	"time"
)

type Member struct {
    MemberID    uint      `gorm:"primaryKey;autoIncrement"`
    FirstName   string    `gorm:"size:50;not null"`
    LastName    string    `gorm:"size:50;not null"`
    Email       string    `gorm:"size:100;uniqueIndex;not null"`
    DateOfBirth time.Time
    Gender      string    `gorm:"size:10"`
    Phone       string    `gorm:"size:15"`
    JoinDate    time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    
    // Relationships
    HealthMetrics    []HealthMetric    `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
    FitnessGoals     []FitnessGoal     `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
    TrainingSessions []TrainingSession `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
    ClassEnrollments []ClassEnrollment `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
}
