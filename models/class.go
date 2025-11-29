package models

import "time"

type Class struct {
    ClassID           uint      `gorm:"primaryKey;autoIncrement"`
    ClassName         string    `gorm:"size:100;not null"`
    TrainerID         *uint     `gorm:"index"` // Pointer for nullable FK
    ScheduleTime      time.Time `gorm:"not null"`
    Duration          int       `gorm:"default:60"`
    Capacity          int       `gorm:"not null;check:capacity > 0"`
    CurrentEnrollment int       `gorm:"default:0;check:current_enrollment >= 0"`
    RoomNumber        string    `gorm:"size:20"`
    
    // Relationships
    Trainer          *Trainer          `gorm:"foreignKey:TrainerID;constraint:OnDelete:SET NULL"`
    ClassEnrollments []ClassEnrollment `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE"`
}