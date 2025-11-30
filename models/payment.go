package models

import "time"

type Payment struct {
    PaymentID   uint      `gorm:"primaryKey;autoIncrement"`
    BillID      uint      `gorm:"not null"`
    Amount      float64   `gorm:"type:decimal(10,2);not null"`
    PaymentDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
    Method      string    `gorm:"size:30"` // cash, credit, debit    
}