package models

import "time"

type ClassEnrollment struct {
    EnrollmentID   uint      `gorm:"primaryKey;autoIncrement"`
    MemberID       uint      `gorm:"not null;index:idx_member_class,unique"`
    ClassID        uint      `gorm:"not null;index:idx_member_class,unique"`
    EnrollmentDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    Status         string    `gorm:"size:20;default:'active'"`
    
    // Relationships
    Member Member `gorm:"foreignKey:MemberID;constraint:OnDelete:CASCADE"`
    Class  Class  `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE"`
}

// Composite unique constraint
func (ClassEnrollment) TableName() string {
    return "class_enrollments"
}