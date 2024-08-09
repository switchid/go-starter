package models

import "time"

type User struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement:true"`
	FirstName string    `gorm:"size:255;not null"`
	LastName  string    `gorm:"size:255;"`
	Email     string    `gorm:"unique;size:100;not null"`
	Password  string    `gorm:"not null"`
	IsActive  bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime:false"`
}

func (User) TableName() string {
	return "tb_users"
}
