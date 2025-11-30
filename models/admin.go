package models

type Admin struct {
	AdminID   uint   `gorm:"primaryKey;autoIncrement"`
	FirstName string `gorm:"size:50;not null"`
	LastName  string `gorm:"size:50;not null"`
	Email     string `gorm:"size:100;uniqueIndex;not null"`
	Phone     string `gorm:"size:15"`
	Role      string `gorm:"size:30;default:'staff'"`
}
