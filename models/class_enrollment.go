package models

import "time"

type ClassEnrollment struct {
	EnrollmentID   uint      `gorm:"primaryKey;autoIncrement"`
	MemberID       uint      `gorm:"not null;index:idx_member_class,unique"`
	ClassID        uint      `gorm:"not null;index:idx_member_class,unique"`
	EnrollmentDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Status         string    `gorm:"size:20;default:'active'"`
}

// Composite unique constraint
func (ClassEnrollment) TableName() string {
	return "class_enrollments"
}
