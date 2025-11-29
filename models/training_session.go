package models

import "time"

type TrainingSession struct {
    SessionID uint      `gorm:"primaryKey;autoIncrement"`
    MemberID  uint      `gorm:"not null;index"`
    TrainerID *uint     `gorm:"index"` // Nullable
    Date      time.Time `gorm:"not null"`
    StartTime time.Time `gorm:"not null"`
    EndTime   time.Time `gorm:"not null"`
    Status    string    `gorm:"size:20;default:'scheduled'"`
    
    // Relationships
    Member  Member   `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
    Trainer *Trainer `gorm:"foreignKey:TrainerID;constraint:OnDelete:SET NULL"`
}