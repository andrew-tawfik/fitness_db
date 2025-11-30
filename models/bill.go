package models

import "time"

type Bill struct {
	BillID    uint      `gorm:"primaryKey;autoIncrement"`
	MemberID  uint      `gorm:"not null"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	IssueDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	DueDate   time.Time
	Status    string `gorm:"size:20;default:'unpaid'"`
}
