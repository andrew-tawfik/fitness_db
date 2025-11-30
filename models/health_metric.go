package models

import "time"

type HealthMetric struct {
	MemberID     uint      `gorm:"primaryKey"` // Part of composite PK
	MetricID     uint      `gorm:"primaryKey"` // Part of composite PK
	RecordedDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Weight       float64   `gorm:"type:decimal(5,2)"`
	Height       float64   `gorm:"type:decimal(5,2)"`
	HeartRate    int
	BodyFatPct   float64 `gorm:"type:decimal(4,2)"`
}

// Composite primary key for weak entity
func (HealthMetric) TableName() string {
	return "health_metrics"
}