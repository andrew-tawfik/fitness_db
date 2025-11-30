package models

import "time"

type FitnessGoal struct {
	MemberID     uint    `gorm:"primaryKey"` // Part of composite PK
	GoalID       uint    `gorm:"primaryKey"` // Part of composite PK
	GoalType     string  `gorm:"size:50"`
	TargetWeight float64 `gorm:"type:decimal(5,2)"`
	TargetDate   time.Time
	Status       string `gorm:"size:20;default:'active'"`
}
