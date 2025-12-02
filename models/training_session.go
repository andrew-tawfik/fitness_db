package models

import "time"

type TrainingSession struct {
	SessionID uint      `gorm:"primaryKey;autoIncrement"`
	MemberID  uint      `gorm:"not null;index"`
	TrainerID *uint     `gorm:"index"` // Nullable
	Date      time.Time `gorm:"type:date;not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	Status    string    `gorm:"size:20;default:'scheduled'"`
}
