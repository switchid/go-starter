package models

import "time"

type TbaUser struct {
	Id            uint      `gorm:"primaryKey;autoIncrement:true"`
	FirstName     string    `gorm:"size:255;not null"`
	LastName      string    `gorm:"size:255;"`
	Address       string    `gorm:"size:255;"`
	Email         string    `gorm:"unique;size:100;not null"`
	Password      string    `gorm:"size:255;not null"`
	UserType      string    `gorm:"size:25;not null"`
	RoleId        uint      `gorm:"not null"`
	IsActive      bool      `gorm:"default:false"`
	RememberToken string    `gorm:"size:255;"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoCreateTime:false"`
}

func (TbaUser) TableName() string {
	return "tba_users"
}
